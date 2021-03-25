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

const composeFile = "docker-compose.yml"

var dockerCompose = sh.RunCmd("docker-compose", "-f", composeFile)

var Default = BuildImage

// BuildImage rebuilds the docker image.
func BuildImage(ctx context.Context) error {
	return dockerCompose("build")
}

// RecreateContainer destroys the container and spins up a new one.
func RecreateContainer(ctx context.Context) {
	mg.CtxDeps(ctx, BuildImage, DestroyContainer)
	mg.CtxDeps(ctx, createContainer)
}

// Bash enters into a new shell inside a running container.
func Bash(ctx context.Context) {
	if !isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS == "windows" {
		sh.RunV("winpty", "docker-compose", "-f", composeFile, "exec", "dev", "bash")
		return
	}
	sh.RunV("docker-compose", "-f", composeFile, "exec", "dev", "bash")
}

func createContainer() error {
	return dockerCompose("-p", "dev-container", "up", "-d")
}

// DestroyContainer destroys the container.
func DestroyContainer() error {
	return dockerCompose("down")
}
