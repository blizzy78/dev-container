//go:build mage

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mattn/go-isatty"
)

const (
	composeFile = "docker-compose.yml"
	projectName = "dev-container"
)

var Default = BuildImage

// BuildImage rebuilds the docker image.
func BuildImage(ctx context.Context) error {
	mg.CtxDeps(ctx, pullGolang, pullUbuntu)

	if err := dockerCompose("build", "--no-cache", "--force-rm"); err != nil {
		return fmt.Errorf("docker compose build: %w", err)
	}

	return nil
}

func pullGolang() error {
	if err := dockerPull("golang"); err != nil {
		return fmt.Errorf("docker pull golang: %w", err)
	}

	return nil
}

func pullUbuntu() error {
	if err := dockerPull("ubuntu"); err != nil {
		return fmt.Errorf("docker pull ubuntu: %w", err)
	}

	return nil
}

// RecreateContainer destroys the container and spins up a new one, optionally recreating the image first.
func RecreateContainer(ctx context.Context, rebuildImage bool) {
	if rebuildImage {
		mg.CtxDeps(ctx, BuildImage, DestroyContainer)
	} else {
		mg.CtxDeps(ctx, DestroyContainer)
	}

	mg.CtxDeps(ctx, CreateContainer)
}

// Bash enters into a new shell inside a running container.
func Bash(ctx context.Context) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		if err := sh.RunV("winpty", "docker", "compose", "-f", composeFile, "-p", projectName, "exec", "dev", "bash"); err != nil {
			return fmt.Errorf("winpty docker compose exec bash: %w", err)
		}

		return nil
	}

	if err := sh.RunV("docker", "compose", "-f", composeFile, "-p", projectName, "exec", "dev", "bash"); err != nil {
		return fmt.Errorf("docker compose exec bash: %w", err)
	}

	return nil
}

// CreateContainer creates the container.
func CreateContainer() error {
	if err := dockerCompose("up", "-d"); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}

	return nil
}

// DestroyContainer destroys the container.
func DestroyContainer() error {
	if err := dockerCompose("down", "--remove-orphans"); err != nil {
		return fmt.Errorf("docker compose down: %w", err)
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
