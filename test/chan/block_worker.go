package main

import (
	"fmt"
	"time"
)

func worker(ch chan struct{}) {
    fmt.Println("Worker is starting...")
    time.Sleep(2 * time.Second)
    fmt.Println("Worker is done!")
    ch <- struct{}{}
}

func main() {
    doneCh := make(chan struct{})
    go worker(doneCh)
    <-doneCh
    fmt.Println("Main function exiting.")
}

