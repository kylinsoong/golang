package main

import (
	"fmt"
)

type WatchedNamespaces struct {
    Namespaces     []string
    NamespaceLabel string
}

func main() {
    watchedNamespaces := WatchedNamespaces{
	Namespaces:     []string{"namespace1", "namespace2"},
	NamespaceLabel: "watched",
    }

    fmt.Println(watchedNamespaces.Namespaces)
    fmt.Println(watchedNamespaces.NamespaceLabel)
}

