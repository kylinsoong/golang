package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
    value int
    mu    sync.Mutex
}

func (c *Counter) increment() {
    c.mu.Lock()
    defer c.mu.Unlock() 
    c.value++
}

func (c *Counter) getValue() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := Counter{}

    for i := 0; i < 5; i++ {
	go func() {
	    for j := 0; j < 3; j++ {
		counter.increment()
		time.Sleep(100 * time.Millisecond)
	    }
	}()
    }

    time.Sleep(2 * time.Second)
    finalValue := counter.getValue()
    fmt.Printf("Final Counter Value: %d\n", finalValue)
}

