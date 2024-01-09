package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 1; i <= 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%d \n", i)
	}
}

func main() {
	go printNumbers()

	for i := 1; i <= 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("A%d \n", i)
	}

}

