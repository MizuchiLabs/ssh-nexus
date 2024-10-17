package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/MizuchiLabs/ssh-nexus/api/client"
	"github.com/MizuchiLabs/ssh-nexus/tools/updater"
)

func main() {
	server := flag.String(
		"server", "127.0.0.1",
		"The address of the server",
	)
	port := flag.Int(
		"port", 8091,
		"The port of the gRPC server",
	)
	version := flag.Bool("version", false, "Show version")
	update := flag.Bool("update", false, "Update to latest version")
	updateCheck := flag.Bool("latest", false, "Check for latest version")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Println(updater.Version)
		return
	}

	if *update || *updateCheck {
		updater.UpdateSelf(updater.Version, *update)
		return
	} else {
		go updater.UpdateSelf(updater.Version, *update)
	}

	addr, err := normalizeURL(*server, *port)
	if err != nil {
		slog.Error("failed to normalize URL", "err", err)
		return
	}

	slog.Info(
		"Starting agent",
		"Version",
		updater.Version,
		"Platform",
		runtime.GOOS+"/"+runtime.GOARCH,
	)
	client.Client(addr)
}

func normalizeURL(input string, port int) (string, error) {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "https://" + input
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s:%d", input, port))
	if err != nil {
		return "", err
	}

	var addr string
	if net.ParseIP(parsedURL.Hostname()) != nil {
		addr = parsedURL.String()
	} else {
		addr = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Hostname())
	}

	return addr, nil
}
