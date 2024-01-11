package main

import (
	"fmt"
	"sync"
)

func numberGenerator(ch chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 1; i <= 5; i++ {
        ch <- i // Send numbers 1 to 5 to the channel
    }
    close(ch) // Close the channel to signal no more data will be sent
}

func squareCalculator(ch chan int, resultCh chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for num := range ch {
	square := num * num
	resultCh <- square // Send squared result to the resultCh channel
    }
    close(resultCh) // Close the resultCh channel to signal no more results will be sent
}

func resultPrinter(resultCh chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for result := range resultCh {
	fmt.Println("Squared Result:", result)
    }
}

func main() {
    numberCh := make(chan int)
    resultCh := make(chan int)
    var wg sync.WaitGroup
    wg.Add(3)
    go numberGenerator(numberCh, &wg)
    go squareCalculator(numberCh, resultCh, &wg)
    go resultPrinter(resultCh, &wg)
    wg.Wait()
}

