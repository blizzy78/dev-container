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
	composeFile         = "docker-compose.yaml"
	composeTemplateFile = "docker-compose.yaml.tmpl"

	projectName = "dev-container"
)

var Default = Build

var composeFileMu sync.Mutex

// Build builds the Docker image.
func Build(ctx context.Context) error {
	return withComposeFile(func() error {
		if err := dockerCompose("pull"); err != nil {
			return fmt.Errorf("docker compose pull: %w", err)
		}

		if err := dockerCompose("build", "--no-cache", "--force-rm"); err != nil {
			return fmt.Errorf("docker compose build: %w", err)
		}

		return nil
	})
}

// Recreate destroys the containers and spins up new ones.
func Recreate(ctx context.Context) {
	mg.SerialCtxDeps(ctx, Destroy, Create)
}

// Zsh enters into a new shell inside a running container.
func Zsh(ctx context.Context) error {
	return withComposeFile(func() error {
		if err := dockerCompose("exec", "vscode", "zsh", "--login"); err != nil {
			return fmt.Errorf("docker compose exec zsh: %w", err)
		}

		return nil
	})
}

// Create creates the containers.
func Create(ctx context.Context) error {
	return withComposeFile(func() error {
		if err := dockerCompose("up", "-d", "--no-start"); err != nil {
			return fmt.Errorf("docker compose up: %w", err)
		}

		if err := dockerCompose("start", "vscode", "pg", "redis"); err != nil {
			return fmt.Errorf("docker compose start: %w", err)
		}

		return nil
	})
}

// Destroy destroys the containers.
func Destroy(ctx context.Context) error {
	return withComposeFile(func() error {
		if err := dockerCompose("down", "--remove-orphans"); err != nil {
			return fmt.Errorf("docker compose down: %w", err)
		}

		return nil
	})
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

func withComposeFile(fn func() error) error {
	removeComposeFile, err := createComposeFile()
	if err != nil {
		return fmt.Errorf("create compose file: %w", err)
	}
	defer removeComposeFile()

	return fn()
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
	}, nil
}
