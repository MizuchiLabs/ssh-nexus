// Package updater contains various utility functions and loads the configuration via envs and for self updating
package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/caarlos0/env/v11"
)

type repository struct {
	URL   string `env:"PB_REPO_URL"   envDefault:"https://github.com/mizuchilabs/ssh-nexus"`
	Owner string `env:"PB_REPO_OWNER" envDefault:"MizuchiLabs"`
	Repo  string `env:"PB_REPO_NAME"  envDefault:"ssh-nexus"`
	Token string `env:"PB_REPO_TOKEN" envDefault:""`
}

type releaseAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	ID          int    `json:"id"`
	Size        int    `json:"size"`
}

type release struct {
	Name      string          `json:"name"`
	Tag       string          `json:"tag_name"`
	Published string          `json:"published_at"`
	URL       string          `json:"html_url"`
	Body      string          `json:"body"`
	Assets    []*releaseAsset `json:"assets"`
	ID        int             `json:"id"`
}

func UpdateSelf(version string, update bool) {
	if IsRunningInDocker() {
		slog.Info("Running in docker, skipping update")
		return
	}

	slog.Info("Fetching release information...")

	latest, err := fetchLatestRelease()
	if err != nil {
		slog.Error("Update failed", "Error", err)
		return
	}

	if !update {
		if compareVersions(
			strings.TrimPrefix(version, "v"),
			strings.TrimPrefix(latest.Tag, "v"),
		) <= 0 {
			slog.Info("You are running the latest version!")
			return
		}
		slog.Info("New version available!", "latest", latest.Tag, "current", version)
		return
	}

	asset := latest.findBinary(filepath.Base(os.Args[0]))
	if asset == nil {
		slog.Info("Unsupported platform", "platform", runtime.GOOS+"/"+runtime.GOARCH)
		return
	}
	exec, err := os.Executable()
	if err != nil {
		slog.Error("Update failed", "Error", err)
		return
	}
	if err := os.Remove(exec); err != nil {
		slog.Error("Failed to remove current executable", "Error", err)
		return
	}

	slog.Info("Downloading...", "release", latest.Tag, "binary", asset.Name)
	if err := downloadFile(asset.DownloadURL, exec); err != nil {
		slog.Error("Failed to download", "Error", err)
		return
	}

	slog.Info("Update success!")
}

// CheckAgent checks if the agent is up to date and downloaded already
func CheckAgent() error {
	_, err := os.Stat(data.AgentDownloadPath)
	if errors.Is(err, os.ErrNotExist) {
		if err = downloadAgent(); err != nil {
			return err
		}
	}

	cmd := exec.Command(data.AgentDownloadPath, "-version")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return err
	}

	latest, err := fetchLatestRelease()
	if err != nil {
		return err
	}
	version := strings.TrimSpace(out.String())
	if compareVersions(
		strings.TrimPrefix(version, "v"),
		strings.TrimPrefix(latest.Tag, "v"),
	) == 1 {
		if err := downloadAgent(); err != nil {
			return err
		}
	}

	return nil
}

func downloadAgent() error {
	latest, err := fetchLatestRelease()
	if err != nil {
		return err
	}
	asset := latest.findBinary("nexus-agent")
	if asset == nil {
		return fmt.Errorf("agent not found")
	}
	slog.Info("Downloading...", "release", latest.Tag, "binary", asset.Name)
	if err := downloadFile(asset.DownloadURL, data.AgentDownloadPath); err != nil {
		return err
	}

	return nil
}

func getRepository() (*repository, error) {
	config := repository{}
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

func fetchLatestRelease() (*release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	repo, err := getRepository()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", repo.URL, repo.Owner, repo.Repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if repo.Token != "" {
		req.Header.Set("Authorization", "token "+repo.Token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("(%d) failed to send latest release request", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &release{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func downloadFile(url string, dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("(%d) failed to send download file request", res.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return err
	}

	if err := out.Chmod(0755); err != nil {
		return err
	}

	return nil
}

func (r *release) findBinary(name string) *releaseAsset {
	var assetName string

	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			assetName = name + "_linux_amd64"
		case "arm64":
			assetName = name + "_linux_arm64"
		case "arm":
			assetName = name + "_linux_armv7"
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			assetName = name + "_darwin_amd64"
		case "arm64":
			assetName = name + "_darwin_arm64"
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			assetName = name + "_windows_amd64"
		case "arm64":
			assetName = name + "_windows_arm64"
		}
	}

	for _, asset := range r.Assets {
		if assetName == asset.Name {
			return asset
		}
	}

	return nil
}

func compareVersions(a, b string) int {
	aSplit := strings.Split(a, ".")
	aTotal := len(aSplit)

	bSplit := strings.Split(b, ".")
	bTotal := len(bSplit)

	limit := aTotal
	if bTotal > aTotal {
		limit = bTotal
	}

	for i := 0; i < limit; i++ {
		var x, y int

		if i < aTotal {
			x, _ = strconv.Atoi(aSplit[i])
		}

		if i < bTotal {
			y, _ = strconv.Atoi(bSplit[i])
		}

		if x < y {
			return 1 // b is newer
		}

		if x > y {
			return -1 // a is newer
		}
	}

	return 0 // equal
}

func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
