package main

import (
    "fmt"
    "strings"
    "os"
    //"time"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

    "github.com/spf13/pflag"
    "golang.org/x/crypto/ssh/terminal"

    log "github.com/kylinsoong/k8s-client-test/pkg/vlogger"
    clog "github.com/kylinsoong/k8s-client-test/pkg/vlogger/console"
)

var (
    flags            *pflag.FlagSet
    globalFlags      *pflag.FlagSet
    kubeFlags        *pflag.FlagSet

    logLevel         *string

    namespaces             *[]string
    inCluster              *bool
    kubeConfig             *string
)

func _init() {
    flags = pflag.NewFlagSet("main", pflag.PanicOnError)
    globalFlags = pflag.NewFlagSet("Global", pflag.PanicOnError)
    kubeFlags = pflag.NewFlagSet("Kubernetes", pflag.PanicOnError)

    var err error
    var width int
    fd := int(os.Stdout.Fd())
    if terminal.IsTerminal(fd) {
        width, _, err = terminal.GetSize(fd)
        if nil != err {
            width = 0
        }
    }

    logLevel = globalFlags.String("log-level", "DEBUG", "Optional, logging level")


    globalFlags.Usage = func() {
        fmt.Fprintf(os.Stderr, "  Global:\n%s\n", globalFlags.FlagUsagesWrapped(width))
    }

    namespaces = kubeFlags.StringArray("namespace", []string{}, "Optional, Kubernetes namespace(s) to watch. If left blank controller will watch all k8s namespaces")
    inCluster = kubeFlags.Bool("running-in-cluster", true, "Optional, if this controller is running in a kubernetes cluster, use the pod secrets for creating a Kubernetes client.")
    kubeConfig = kubeFlags.String("kubeconfig", "./config", "Optional, absolute path to the kubeconfig file")
  
    kubeFlags.Usage = func() {
        fmt.Fprintf(os.Stderr, "  Kubernetes:\n%s\n", kubeFlags.FlagUsagesWrapped(width))
    }

    flags.AddFlagSet(globalFlags)
    flags.AddFlagSet(kubeFlags)

    flags.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
            globalFlags.Usage()
            kubeFlags.Usage()
    }
}

func initLogger(logLevel string) error {
    log.RegisterLogger(log.LL_MIN_LEVEL, log.LL_MAX_LEVEL, clog.NewConsoleLogger())

    if ll := log.NewLogLevel(logLevel); nil != ll {
        log.SetLogLevel(*ll)
    } else {
        return fmt.Errorf("Unknown log level requested: %s\n" + "    Valid log levels are: DEBUG, INFO, WARNING, ERROR, CRITICAL", logLevel)
    }
    return nil
}

func init() {
    _init()
}

func getKubeConfig() (*rest.Config, error) {
    var config *rest.Config
    var err error
    if *inCluster {
        config, err = rest.InClusterConfig()
    } else {
        config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
    }
        
    if err != nil {
        log.Fatalf("[INIT] error creating configuration: %v", err)
        return nil, err
    }
        
    return config, nil
}

func main() {

    *logLevel = strings.ToUpper(*logLevel)
    initLogger(*logLevel)

    log.Infof("[INIT] Starting: K8S CLIENT TEST")

    log.Infof("[INIT] logLevel: %s, inCluster: %t, kubeConfig: %s", *logLevel, *inCluster, *kubeConfig)

    config, err := getKubeConfig()
    if err != nil {
        os.Exit(1)
    }

    kubeClient, err = kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("[INIT] error connecting to the client: %v", err)
        os.Exit(1)
    }


    
}
