package main

import "fmt"

func getEndpoints(selector, namespace, svcNamespace string) []string {
    endpoints := []string{
        fmt.Sprintf("%s", selector),
        fmt.Sprintf("%s", namespace),
        fmt.Sprintf("%s", svcNamespace),
    }
    return endpoints
}

func main() {
    selector := "my-service"
    namespace := "default"
    svcNamespace := "prod"

    endpoints := getEndpoints(selector, namespace, svcNamespace)
    fmt.Println(endpoints)

    endpoints = getEndpoints(selector, namespace)
    fmt.Println(endpoints)
}

