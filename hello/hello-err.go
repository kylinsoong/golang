package main

import (
    "fmt"
    "cloudadc.github.io/greetings"
)

func main() {
    fmt.Println(greetings.Hello("Kylin"))
    fmt.Println(greetings.Hello(""))
}
