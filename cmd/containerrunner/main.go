package main

import (
	"context"
	"os/signal"
	"syscall"
)

func main() {
	run()
}

func run() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
}
