package main

import (
    "fmt"
    "math/rand"
)

var (
    globalMap map[string]bool
)

func generateIP() string{
    ip := fmt.Sprintf("10.244.%d.%d", rand.Intn(256), rand.Intn(256))
    if globalMap[ip] {
        return generateIP()
    } else {
        globalMap[ip] = true
        return ip
    }

}

func main() {
    globalMap = make(map[string]bool)
    for i := 0; i < 10; i++ {
        fmt.Println(generateIP(), globalMap)
    }
}

