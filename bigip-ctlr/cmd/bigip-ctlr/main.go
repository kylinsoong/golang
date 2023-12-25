package main

import (
    "fmt"
    "strings"
    "os"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

    cisAgent "github.com/kylinsoong/bigip-ctlr/pkg/agent"
    "github.com/kylinsoong/bigip-ctlr/pkg/appmanager"
    "github.com/kylinsoong/bigip-ctlr/pkg/resource"
)

type globalSection struct {
        LogLevel       string `json:"log-level,omitempty"`
        VerifyInterval int    `json:"verify-interval,omitempty"`
        VXLANPartition string `json:"vxlan-partition,omitempty"`
        DisableLTM     bool   `json:"disable-ltm,omitempty"`
        DisableARP     bool   `json:"disable-arp,omitempty"`
}

type bigIPSection struct {
        BigIPUsername   string   `json:"username,omitempty"`
        BigIPPassword   string   `json:"password,omitempty"`
        BigIPURL        string   `json:"url,omitempty"`
        BigIPPartitions []string `json:"partitions,omitempty"`
}

var (
    kubeClient         kubernetes.Interface
)

func getKubeConfig() (*rest.Config, error) {

    var config *rest.Config

    config = &rest.Config{
        Host:            "https://dummy-kube-api-server",
	BearerToken:     "dummy-bearer-token",
	TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	Timeout:         30, // Set the timeout in seconds
    }

    return config, nil
}

func main() {

    fmt.Println("F5 BIG-IP Controller Start")

    dgPath := strings.Join([]string{resource.DEFAULT_PARTITION, "Shared"}, "/")
    fmt.Printf("BIG-IP Data Group Path: %s\n", dgPath)

    appmanager.RegisterBigIPSchemaTypes()

    config, err := getKubeConfig()
    if err != nil {
        os.Exit(1)
    }

    kubeClient, err = kubernetes.NewForConfig(config)
    if err != nil {
        os.Exit(1)
    }

    fmt.Printf("create kubeClient via %s\n", config.Host)

    agent := "as3"
    disableLTM := false
    if agent == cisAgent.AS3Agent {
            disableLTM = true
    }

    disableARP := false

    gs := globalSection{
        LogLevel:       "DEBUG",
        VerifyInterval: 60,
        VXLANPartition: "Common",
        DisableLTM:     disableLTM,
        DisableARP:     disableARP,
    }

    bs := bigIPSection{
        BigIPUsername:   "admin",
        BigIPPassword:   "admin",
        BigIPURL:        "https://10.1.10.240",
        BigIPPartitions: []string{"partition1", "partition2"},
    }

    fmt.Printf("Create Global Section and BIG-IP Section struct\n%+v\n", gs)
    fmt.Printf("%v\n", bs)
}
