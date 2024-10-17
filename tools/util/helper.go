// Package util provides common utility functions
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		if value != "" {
			return value
		}
	}
	return fallback
}

// Execute a goroutine of a task
func Execute(fn func()) {
	go func() {
		defer recoverPanic()
		fn()
	}()
}

// Write the error to console when a goroutine of a task panicking.
func recoverPanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		fmt.Println(err)
	}
}

func IsIP(host string) bool {
	u, err := url.Parse(host)
	if err == nil && u.Host != "" {
		host = u.Host
	}
	host = strings.Split(host, ":")[0] // Remove port if present

	// Check if it's a valid IP address
	return net.ParseIP(host) != nil
}

func IsValidDomain(domain string) bool {
	if IsIP(domain) {
		return false
	}

	// A very basic regex for domain validation
	re := regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	return re.MatchString(domain)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		slog.Error("failed to get outbound IP", "err", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetPublicIP() net.IP {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil
	}

	var data struct {
		Query string `json:"query"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil
	}

	remoteAddr := net.ParseIP(data.Query)

	return remoteAddr
}

// Diff returns the elements in a that are not in b
func Diff(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func GetLeaseDuration(ttl interface{}, defaultTTL, maxTTL int) time.Duration {
	var leaseDuration time.Duration
	if ttl != nil {
		TTL, _ := strconv.Atoi(ttl.(string))
		leaseDuration = time.Duration(TTL) * time.Second
	} else {
		leaseDuration = time.Duration(defaultTTL) * time.Second
	}
	if leaseDuration > time.Duration(maxTTL)*time.Second {
		leaseDuration = time.Duration(maxTTL) * time.Second
	}
	return leaseDuration
}
