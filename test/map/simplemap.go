package main

import "fmt"

func main() {

    processedResources := make(map[string]bool)

    processedResources["foo.yaml"] = true
    processedResources["bar.yaml"] = false
    processedResources["zoo.yaml"] = false

    for key, value := range processedResources {
        fmt.Printf("%s: %v\n", key, value)
    }

    fmt.Println(processedResources["zoo.yaml"])

    value, exists := processedResources["coo.yaml"]
    if exists {
        fmt.Printf("coo.yaml: %v\n", value) 
    } else {
        fmt.Println("coo.yaml not exist")
    }
}

