//go:build mage

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/magefile/mage/mg"
)

type httpError string

func bashStdin(r io.Reader, args ...string) error {
	c := exec.Command("bash", args...)
	c.Stdin = r
	o, err := c.CombinedOutput()
	if mg.Verbose() {
		_, _ = os.Stdout.Write(o)
	}

	if err != nil {
		return fmt.Errorf("run bash script: %w", err)
	}

	return nil
}

func mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func copyToFile(r io.Reader, path string, perm fs.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, r); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

func downloadAs(ctx context.Context, url string, path string) error {
	b, err := download(ctx, url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func download(ctx context.Context, url string) ([]byte, error) {
	c := http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,

			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, httpError(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return body, nil
}

func (e httpError) Error() string {
	return "HTTP error: " + string(e)
}
