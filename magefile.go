//go:build mage

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mattn/go-isatty"
)

const (
	composeFile         = "docker-compose.yml"
	composeTemplateFile = "docker-compose.tmpl.yml"

	projectName = "dev-container"
)

var (
	Default = BuildImage

	composeFileMu sync.Mutex
)

// BuildImage rebuilds the docker image.
func BuildImage(ctx context.Context) error {
	mg.CtxDeps(ctx, pullGolang, pullArchLinux)

	removeComposeFile, err := createComposeFile()
	if err != nil {
		return fmt.Errorf("create compose file: %w", err)
	}
	defer removeComposeFile()

	if err = dockerCompose("build", "--no-cache", "--force-rm"); err != nil {
		return fmt.Errorf("docker compose build: %w", err)
	}

	return nil
}

// RecreateContainer destroys the container and spins up a new one, optionally recreating the image first.
func RecreateContainer(ctx context.Context, rebuildImage bool) {
	if rebuildImage {
		mg.CtxDeps(ctx, DestroyContainer, BuildImage)
	} else {
		mg.CtxDeps(ctx, DestroyContainer)
	}

	mg.CtxDeps(ctx, CreateContainer)
}

// Bash enters into a new shell inside a running container.
func Bash(ctx context.Context) error {
	removeComposeFile, err := createComposeFile()
	if err != nil {
		return fmt.Errorf("create compose file: %w", err)
	}
	defer removeComposeFile()

	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		if err := sh.RunV("winpty", "docker", "compose", "-f", composeFile, "-p", projectName, "exec", "dev", "bash", "--login"); err != nil {
			return fmt.Errorf("winpty docker compose exec bash: %w", err)
		}

		return nil
	}

	if err := sh.RunV("docker", "compose", "-f", composeFile, "-p", projectName, "exec", "dev", "bash", "--login"); err != nil {
		return fmt.Errorf("docker compose exec bash: %w", err)
	}

	return nil
}

// CreateContainer creates the container.
func CreateContainer(ctx context.Context) error {
	removeComposeFile, err := createComposeFile()
	if err != nil {
		return fmt.Errorf("create compose file: %w", err)
	}
	defer removeComposeFile()

	if err := dockerCompose("up", "-d"); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}

	return nil
}

// DestroyContainer destroys the container.
func DestroyContainer(ctx context.Context) error {
	removeComposeFile, err := createComposeFile()
	if err != nil {
		return fmt.Errorf("create compose file: %w", err)
	}
	defer removeComposeFile()

	if err := dockerCompose("down", "--remove-orphans"); err != nil {
		return fmt.Errorf("docker compose down: %w", err)
	}

	return nil
}

func pullGolang() error {
	if err := dockerPull("golang"); err != nil {
		return fmt.Errorf("docker pull golang: %w", err)
	}

	return nil
}

func pullArchLinux() error {
	if err := dockerPull("archlinux"); err != nil {
		return fmt.Errorf("docker pull archlinux: %w", err)
	}

	return nil
}

func dockerPull(args ...string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		if err := sh.RunV("winpty", append([]string{"docker", "pull"}, args...)...); err != nil {
			return fmt.Errorf("winpty docker pull: %w", err)
		}
	}

	if err := sh.RunV("docker", append([]string{"pull"}, args...)...); err != nil {
		return fmt.Errorf("docker pull: %w", err)
	}

	return nil
}

func dockerCompose(args ...string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		if err := sh.RunV("winpty", append([]string{"docker", "compose", "-f", composeFile, "-p", projectName}, args...)...); err != nil {
			return fmt.Errorf("winpty docker compose: %w", err)
		}

		return nil
	}

	if err := sh.RunV("docker", append([]string{"compose", "-f", composeFile, "-p", projectName}, args...)...); err != nil {
		return fmt.Errorf("docker compose: %w", err)
	}

	return nil
}

func createComposeFile() (func(), error) {
	composeFileMu.Lock()

	tmp, err := template.ParseFiles(composeTemplateFile)
	if err != nil {
		return nil, fmt.Errorf("parse compose template: %w", err)
	}

	file, err := os.OpenFile(composeFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", composeFile, err)
	}

	fileClosed := false
	defer func() {
		if !fileClosed {
			_ = file.Close()
		}
	}()

	data := struct {
		GOOS string
	}{
		GOOS: runtime.GOOS,
	}

	if err := tmp.Execute(file, data); err != nil {
		return nil, fmt.Errorf("execute compose template: %w", err)
	}

	fileClosed = true
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close %s: %w", composeFile, err)
	}

	return func() {
		defer composeFileMu.Unlock()
		_ = os.Remove(composeFile)
	}, nil
}
