package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"time"
)

var (
	envDelay     = "KEEP_DELAY"
	defaultDelay = 50 * time.Millisecond
	delay        = defaultDelay
)

func main() {
	envDelay := os.Getenv(envDelay)
	if envDelay != "" {
		parsed, err := time.ParseDuration(envDelay)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("parsing delay: %w", err))
			printUsageAndExit()
		}
		delay = parsed
	}

	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	name := os.Args[1]
	args := os.Args[2:]

	done := make(chan struct{})
	go func() {
		defer func() {
			close(done)
		}()

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

		for restartCount := 0; true; restartCount++ {
			// run command
			cmd := exec.CommandContext(ctx, name, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run() // ignore error

			// break loop if context was canceled
			canceled := ctx.Err() != nil
			printCommandExited(restartCount, cmd.ProcessState.ExitCode(), !canceled)
			if canceled {
				break
			}

			// execute delay
			ctxsleep(ctx, delay)

			// check context once again
			canceled = ctx.Err() != nil
			if canceled {
				break
			}
		}
	}()

	<-done
}

func printUsageAndExit() {
	kcmd := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, `usage: %s [command --args]
where "command --args" is a long running operation that will be restarted by %s when it exits.

Configuration options:
delay | env KEEP_DELAY | wait time between restarts (default: %s) (actual: %s)
`, kcmd, kcmd, defaultDelay, delay)
	os.Exit(1)
}

func printCommandExited(restartCount, code int, willRestart bool) {
	if willRestart {
		fmt.Fprintf(os.Stderr, "[restarts: %d] Command exited with error code: %d. Restarting after delay of %s.\n", restartCount, code, delay)
	} else {
		fmt.Fprintf(os.Stderr, "[restarts: %d] Command exited with error code: %d.\n", restartCount, code)
	}
}

func ctxsleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer func() {
		if t.Stop() {
			return
		}
		go func() {
			<-t.C
		}()
	}()

	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
