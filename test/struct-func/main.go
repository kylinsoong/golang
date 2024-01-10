package main

import (
	"fmt"
)

type Manager struct {
    queueLen            int
    processAgentLabels  func(map[string]string, string, string) bool
}

func customProcessAgentLabels(labels map[string]string, namespace string, name string) bool {
    fmt.Printf("Custom Processing Agent Labels: %v, Namespace: %s, Name: %s\n", labels, namespace, name)
    return true 
}

func main() {
    appMgr := Manager{
        queueLen:           10,
        processAgentLabels: customProcessAgentLabels,
    }
    appMgr.processAgentLabels(map[string]string{"key": "value"}, "exampleNamespace", "exampleName")
}

