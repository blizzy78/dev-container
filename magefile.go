//go:build mage
// +build mage

package main

import (
	"context"
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
	return dockerCompose("build")
}

func pullGolang() error {
	return dockerPull("golang")
}

func pullUbuntu() error {
	return dockerPull("ubuntu")
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
func Bash(ctx context.Context) {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		sh.RunV("winpty", "docker-compose", "-f", composeFile, "-p", projectName, "exec", "dev", "bash")
		return
	}
	sh.RunV("docker-compose", "-f", composeFile, "-p", projectName, "exec", "dev", "bash")
}

// CreateContainer creates the container.
func CreateContainer() error {
	return dockerCompose("up", "-d")
}

// DestroyContainer destroys the container.
func DestroyContainer() error {
	return dockerCompose("down")
}

func dockerPull(args ...string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		return sh.RunV("winpty", append([]string{"docker", "pull"}, args...)...)
	}
	return sh.RunV("docker", append([]string{"pull"}, args...)...)
}

func dockerCompose(args ...string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		return sh.RunV("winpty", append([]string{"docker", "compose", "-f", composeFile, "-p", projectName}, args...)...)
	}
	return sh.RunV("docker", append([]string{"compose", "-f", composeFile, "-p", projectName}, args...)...)
}
