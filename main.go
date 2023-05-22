package main

import (
	"github.com/bthuilot/dockerleaks/cmd"
	"os"
	"os/signal"
)

func main() {
	registerInterrupt()

	cmd.Execute()
}

func registerInterrupt() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		<-stopChan
		panic("exiting.")
	}()
}
