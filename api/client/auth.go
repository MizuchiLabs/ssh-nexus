package client

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"golang.org/x/net/http2"
)

func GetToken() ([]byte, error) {
	token, err := os.ReadFile(data.Token)
	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		return nil, fmt.Errorf("token not found")
	}

	return bytes.TrimSpace(token), nil
}

func GetAgentID() ([]byte, error) {
	machineID, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return nil, err
	}

	if len(machineID) > 0 {
		return bytes.TrimSpace(machineID), nil
	}

	return createAgentID()
}

func LoadCredentials(addr string) (*http.Client, error) {
	// Since we're using self signed certificates, we need to skip the verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	serverCa, err := client.Get(fmt.Sprintf("%s/ca.crt", addr))
	if err != nil {
		return nil, err
	}
	defer serverCa.Body.Close()

	if serverCa.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get server ca: %d", serverCa.StatusCode)
	}

	serverCaBytes, err := io.ReadAll(serverCa.Body)
	if err != nil {
		return nil, err
	}

	// Create a certificate pool and append the CA certificate
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(serverCaBytes); !ok {
		return nil, err
	}

	return &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caPool,
				MinVersion: tls.VersionTLS12,
			},
		},
	}, nil
}

func createAgentID() ([]byte, error) {
	machineID := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, machineID)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(data.Path("agent-id"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() == 0 {
		// The file is newly created, write some content
		if _, err := file.WriteString(string(machineID)); err != nil {
			return nil, err
		}
	} else {
		// The file already exists, read the content
		if _, err := io.ReadFull(file, machineID); err != nil {
			return nil, err
		}
	}

	return bytes.TrimSpace(machineID), nil
}
