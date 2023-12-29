package main

import (
    "fmt"
    "sync"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProcessedHostPath struct {
    sync.Mutex
    processedHostPathMap map[string]metav1.Time
}

type Manager struct {
    processedResources  map[string]bool
    processedHostPath   ProcessedHostPath
    nplStore            map[string]NPLAnnoations
}

type NPLAnnotation struct {
        PodPort  int32  `json:"podPort"`
        NodeIP   string `json:"nodeIP"`
        NodePort int32  `json:"nodePort"`
}

type NPLAnnoations []NPLAnnotation

func main() {

    manager := Manager{}

    manager.processedResources = make(map[string]bool)
    manager.processedHostPath.processedHostPathMap = make(map[string]metav1.Time)
    manager.nplStore = make(map[string]NPLAnnoations)

    manager.processedResources["a.yaml"] = true
    manager.processedResources["b.yaml"] = false
    manager.processedResources["c.yaml"] = false
    fmt.Printf("%+v\n", manager)

    manager.processedHostPath.processedHostPathMap["test1"] = metav1.Now()
    manager.processedHostPath.processedHostPathMap["test2"] = metav1.Now()
    manager.processedHostPath.processedHostPathMap["test3"] = metav1.Now()
    fmt.Printf("%+v\n", manager)

    annotations := NPLAnnoations{
		{PodPort: 80, NodeIP: "192.168.1.1", NodePort: 30001},
		{PodPort: 8080, NodeIP: "192.168.1.2", NodePort: 30002},
    }
    manager.nplStore["zdhyw"] = annotations

    fmt.Printf("%+v\n", manager)
    fmt.Printf("%+v\n", manager.processedResources)
    fmt.Printf("%+v\n", manager.processedHostPath)
    fmt.Printf("%+v\n", manager.nplStore)
    
}

