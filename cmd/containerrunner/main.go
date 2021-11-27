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

	for {
		done, err := runSignals(context.Background(), configPath)
		if err != nil {
			panic(err)
		}

		if done {
			break
		}
	}
}

func runSignals(ctx context.Context, configPath string) (bool, error) {
	fmt.Printf("load configuration from file: %s\n", configPath)

	config, loaded, err := loadConfigDefault(configPath)
	if err != nil {
		return true, fmt.Errorf("load configuration: %w", err)
	}

	if !loaded {
		fmt.Fprintf(os.Stderr, "%s: file not found, ignoring\n", configPath)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	hup := make(chan os.Signal, 1)
	defer close(hup)

	signal.Notify(hup, syscall.SIGHUP)
	defer signal.Stop(hup)

	hupDone := make(chan struct{})
	defer close(hupDone)

	huped := false

	go func() {
		select {
		case <-hupDone:

		case <-hup:
			fmt.Println("received SIGHUP")

			huped = true

			cancel()
		}
	}()

	if err := run(ctx, config); err != nil {
		return true, err
	}

	if huped {
		return false, nil
	}

	return true, nil
}

func run(ctx context.Context, config *configuration) error {
	cronConfig := &cronConfiguration{}
	if config.Cron != nil {
		cronConfig = config.Cron
	}

	stop := startCron(cronConfig)
	defer stop()

	<-ctx.Done()

	return nil
}

func startCron(config *cronConfiguration) func() {
	sched := gocron.NewScheduler(time.Local)

	now := time.Now()

	for _, job := range config.Jobs {
		if err := scheduleJob(job, sched, now); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("schedule job '%s': %w", job.Command, err).Error())
		}
	}

	sched.StartAsync()

	return sched.Stop
}

func scheduleJob(job *cronJobConfiguration, sched *gocron.Scheduler, now time.Time) error {
	sched.Every(job.Every)

	if job.Delay != "" {
		delay, err := time.ParseDuration(job.Delay)
		if err != nil {
			return fmt.Errorf("delay '%s': %w", job.Command, job.Delay, err)
		}

		sched.StartAt(now.Add(delay))
	}

	if _, err := sched.Do(runJob, job); err != nil {
		return fmt.Errorf("schedule: %w", job.Command, err)
	}

	return nil
}

func runJob(job *cronJobConfiguration) {
	fmt.Printf("run job: %s\n", job.Command)

	cmd := exec.CommandContext(context.Background(), job.Command, job.Args...) //nolint:gosec // command is user input

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("run job '%s': %w", job.Command, err).Error())
		return
	}

	fmt.Print(string(output))
}
