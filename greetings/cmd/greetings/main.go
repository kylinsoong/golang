package main

import (
    "fmt"
    "log"

    "rsc.io/quote"
    "github.com/kylinsoong/golang/greetings/pkg/greetings"
)

func main() {
    fmt.Println("Hello, World!")

    fmt.Println(quote.Glass())
    fmt.Println(quote.Go())
    fmt.Println(quote.Hello())
    fmt.Println(quote.Opt())

    log.SetPrefix("greetings: ")
    log.SetFlags(0)

    names := []string{"Gladys", "Samantha", "Darrin", "Kylin"}
    messages, err := greetings.Hellos(names)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(messages)

}

