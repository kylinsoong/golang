package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
    fmt.Println("Started to run tasks...")
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
    sig := <-signals
    fmt.Printf("Received signal: %v\n", sig)
}

