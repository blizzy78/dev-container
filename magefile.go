//+build mage

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

var dockerCompose = sh.RunCmd("docker-compose", "-f", composeFile, "-p", projectName)
var dockerPull = sh.RunCmd("docker", "pull")

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

// RecreateContainer destroys the container and spins up a new one.
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

func CreateContainer() error {
	return dockerCompose("up", "-d")
}

// DestroyContainer destroys the container.
func DestroyContainer() error {
	return dockerCompose("down")
}
