package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(stopCh <-chan struct{}) {
    fmt.Printf("run %v\n", stopCh)
}

func main() {
    stopCh := make(chan struct{})
    Run(stopCh)
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    sig := <-sigs
    close(stopCh)
    fmt.Printf("Exiting - signal: %v\n", sig)
}

