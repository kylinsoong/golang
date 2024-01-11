package main

import "fmt"

func main() {
    defer fmt.Println("This will be executed third.")
    defer fmt.Println("This will be executed second.")
    defer fmt.Println("This will be executed first.")
    fmt.Println("Hello, Go!")
}

