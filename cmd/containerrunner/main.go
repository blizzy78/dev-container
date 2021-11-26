package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
)

var errMissingConfigPath = errors.New("missing configuration file path")

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", configPath, "path to config file")

	flag.Parse()

	if configPath == "" {
		panic(errMissingConfigPath)
	}

	config, err := loadConfig(configPath)
	if err != nil {
		err = fmt.Errorf("load configuration: %w", err)

		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		fmt.Fprintln(os.Stderr, fmt.Errorf("ignoring error: %w", err).Error())

		config = &configuration{}
	}

	if err = run(context.Background(), config); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, config *configuration) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cronConfig := &cronConfiguration{}
	if config.Cron != nil {
		cronConfig = config.Cron
	}

	stopCron, err := startCron(cronConfig)
	if err != nil {
		return fmt.Errorf("start cron: %w", err)
	}

	defer stopCron()

	<-ctx.Done()

	return nil
}

func startCron(config *cronConfiguration) (func(), error) {
	sched := gocron.NewScheduler(time.Local)

	now := time.Now()

	for _, job := range config.Jobs {
		sched.Every(job.Every)

		if job.Delay != "" {
			delay, err := time.ParseDuration(job.Delay)
			if err != nil {
				return nil, fmt.Errorf("job '%s': initial delay '%s': %w", job.Command, job.Delay, err)
			}

			sched.StartAt(now.Add(delay))
		}

		if _, err := sched.Do(runJob, job); err != nil {
			return nil, fmt.Errorf("schedule job '%s': %w", job.Command, err)
		}
	}

	sched.StartAsync()

	return func() {
		sched.Stop()
	}, nil
}

func runJob(job *cronJobConfiguration) {
	cmd := exec.CommandContext(context.Background(), job.Command, job.Args...) //nolint:gosec // command is user input

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("run job '%s': %w", job.Command, err).Error())
		return
	}

	fmt.Print(string(output))
}
