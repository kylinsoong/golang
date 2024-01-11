package main

import (
	"fmt"
	"time"
	"k8s.io/apimachinery/pkg/util/wait"
)

func exampleWork() {
    fmt.Println("Doing some work...")
    time.Sleep(2 * time.Second)
}

func main() {
    stopCh := make(chan struct{})
    go wait.Until(exampleWork, time.Second, stopCh)
    time.Sleep(5 * time.Second)
    close(stopCh)
    time.Sleep(1 * time.Second)
    fmt.Println("Main goroutine exiting...")
}

