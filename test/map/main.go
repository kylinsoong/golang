package main

import "fmt"

func main() {

    m := map[string]int{
        "apple":  1,
        "banana": 2,
        "orange": 3,
    }

    fmt.Println("apple:", m["apple"]) 
    fmt.Println("banana:", m["banana"])

    m["banana"] = 4
    fmt.Println("banana:", m["banana"])

    m["grape"] = 10
    
    fmt.Println(m)

    processedResources := make(map[string]bool)

    processedResources["file1.txt"] = true
    processedResources["file2.txt"] = false

    fmt.Println("file1.txt processed:", processedResources["file1.txt"]) // true
    fmt.Println("file2.txt processed:", processedResources["file2.txt"])
}

