/*
Copyright (C) 2025 erdii <me@erdii.engineering>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

func main() {
	sleepDuration := flag.Int("sleep-duration", 15, "amount of seconds to sleep between loop checks")
	timeout := flag.Int("timeout", 600, "amount of seconds the queue must be empty before webhook will be called")
	webhook := flag.String("webhook", os.Getenv("WEBHOOK"), "webhook address to be called with queue length as payload. Env: WEBHOOK")
	flag.Parse()

	if *webhook == "" {
		slog.Error("webhook flag cannot be empty")
		flag.PrintDefaults()
		os.Exit(1)
	}

	_, err := url.Parse(*webhook)
	if err != nil {
		slog.Error("webhook url was not parsed correctly",
			"error", err)
		os.Exit(1)
	}

	for {
		n, err := getQueueLength()
		if err != nil {
			slog.Error("getting print queue length", "error", err)
			os.Exit(1)
		}

		if n == 0 {
			slog.Debug("queue empty")
			time.Sleep(time.Duration(*sleepDuration) * time.Second)
			continue
		}

		slog.Info("turn on - send queue length")
		if err := sendQueueLength(*webhook, n); err != nil {
			slog.Error("sending queue length",
				"error", err)
		}

		iterationCount := *timeout / *sleepDuration
		for i := 0; i < iterationCount; i++ {
			slog.Info("queue not empty",
				"length", n,
				"iteration", i+1)
			time.Sleep(time.Duration(*sleepDuration) * time.Second)

			n, err = getQueueLength()
			if err != nil {
				slog.Error("getting print queue length", "error", err)
				os.Exit(1)
			}
			if n != 0 {
				i = -1
			}
		}

		slog.Info("turn off - send queue length")
		if err := sendQueueLength(*webhook, n); err != nil {
			slog.Error("sending queue length",
				"error", err)
		}
	}
}

func getQueueLength() (int, error) {
	cmd := exec.Command("lpstat", "-o", "-W", "not-completed")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("executing lpstat: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var n int
	for scanner.Scan() {
		n++
	}
	if scanner.Err() != nil {
		return 0, fmt.Errorf("scanning lpstat output lines: %w", err)
	}
	return n, nil
}

func sendQueueLength(webhook string, length int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook+fmt.Sprintf("?length=%d", length), bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return err
}
