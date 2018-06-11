package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"facua.org/compartidosd/client/app"
	"facua.org/compartidosd/common"
)

var mounted []common.IndexEntry

// This daemon works by periodiaclly executing ticks (see the doc for app.Tick)
// This function will periodically execute the ticks until SIGINT is received,
// and will make sure to wait for a tick to finish should de signal be received
// during one, and to cleanup before exiting.
func main() {
	var tickIntervalMs int
	tickIntervalEnv := os.Getenv("TICK_INTERVAL_MS")

	if tickIntervalEnv == "" {
		tickIntervalMs = 30000
	} else {
		var err error
		tickIntervalMs, err = strconv.Atoi(tickIntervalEnv)

		if err != nil {
			fmt.Println("Failed to parse env TICK_INTERVAL_MS")
			panic(err)
		}
	}

	tickInterval := time.Duration(tickIntervalMs) * time.Millisecond

	appRunning := true

	sigIntChannel := make(chan os.Signal, 1)
	stopWaitingChannel := make(chan bool, 1)
	signal.Notify(sigIntChannel, os.Interrupt)

	go func() {
		for range sigIntChannel {
			appRunning = false
			stopWaitingChannel <- true
		}
	}()

	app.Start()

	for appRunning {
		// Wait the tick interval, unless we receive SIGINT
		select {
		case <-stopWaitingChannel:
		case <-time.After(tickInterval):
		}

		// The app might have received SIGINT during the sleep
		// In that case, don't tick
		if appRunning {
			app.Tick()
		}
	}

	app.Stop()
}
