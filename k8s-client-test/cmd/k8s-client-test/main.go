package main

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "os"
    "time"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"

    "github.com/spf13/pflag"
    "golang.org/x/crypto/ssh/terminal"

    "github.com/kylinsoong/k8s-client-test/pkg/appmanager"
    "github.com/kylinsoong/k8s-client-test/pkg/pollers" 
    log "github.com/kylinsoong/k8s-client-test/pkg/vlogger"
    clog "github.com/kylinsoong/k8s-client-test/pkg/vlogger/console"
)

const (
    versionPathk8s         = "/version"
)

var (
    flags            *pflag.FlagSet
    globalFlags      *pflag.FlagSet
    kubeFlags        *pflag.FlagSet

    logLevel         *string
    nodePollInterval *int

    namespaceLabel         *string
    namespaces             *[]string
    inCluster              *bool
    kubeConfig             *string
    manageConfigMaps       *bool
    manageIngress          *bool
    hubMode                *bool
    nodeLabelSelector      *string

    kubeClient         kubernetes.Interface
    agRspChan          chan interface{}
    eventChan          chan interface{}
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
    nodePollInterval = globalFlags.Int("node-poll-interval", 30, "Optional, interval (in seconds) at which to poll for cluster nodes.")

    globalFlags.Usage = func() {
        fmt.Fprintf(os.Stderr, "  Global:\n%s\n", globalFlags.FlagUsagesWrapped(width))
    }

    namespaceLabel = kubeFlags.String("namespace-label", "", "Optional, used to watch for namespaces with this label")
    namespaces = kubeFlags.StringArray("namespace", []string{}, "Optional, Kubernetes namespace(s) to watch. If left blank controller will watch all k8s namespaces")
    inCluster = kubeFlags.Bool("running-in-cluster", true, "Optional, if this controller is running in a kubernetes cluster, use the pod secrets for creating a Kubernetes client.")
    kubeConfig = kubeFlags.String("kubeconfig", "./config", "Optional, absolute path to the kubeconfig file")
    manageIngress = kubeFlags.Bool("manage-ingress", true, "Optional, specify whether or not to manage Ingress resources")
    manageConfigMaps = kubeFlags.Bool("manage-configmaps", true, "Optional, specify whether or not to manage ConfigMap resources")
    hubMode = kubeFlags.Bool("hubmode", false, "Optional, specify whether or not to manage ConfigMap resources in hub-mode")
    nodeLabelSelector = kubeFlags.String("node-label-selector", "", "Optional, used to watch only for nodes with this label")

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

func getk8sVersion() string {
    var versionInfo map[string]string
    var err error      
    var vInfo []byte   
    rc := kubeClient.Discovery().RESTClient()
    if vInfo, err = rc.Get().AbsPath(versionPathk8s).DoRaw(context.TODO()); err == nil  {
        if er := json.Unmarshal(vInfo, &versionInfo); er == nil  {
           return fmt.Sprintf(versionInfo["gitVersion"])
        }
    }
    return ""
}

func getAppManagerParams() appmanager.Params {
    return appmanager.Params{
        ManageConfigMaps:       *manageConfigMaps,
        ManageIngress:          *manageIngress,
        AgRspChan:              agRspChan,
    }
}

func GetNamespaces(appMgr *appmanager.Manager) {
    if len(*namespaces) != 0 && len(*namespaceLabel) == 0 {
        appMgr.WatchedNS.Namespaces = *namespaces
        log.Infof("[INIT] watched namespace: %v", appMgr.WatchedNS.Namespaces)
    }
    if len(*namespaces) == 0 && len(*namespaceLabel) != 0 {
        appMgr.WatchedNS.NamespaceLabel = *namespaceLabel
    }
}

func setupNodePolling(appMgr *appmanager.Manager, np pollers.Poller, eventChanl <-chan interface{}, kubeClient kubernetes.Interface,) error { 

    return nil
}

func main() {

    err := flags.Parse(os.Args)
    if nil != err {
        os.Exit(1)
    }

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

    agRspChan = make(chan interface{}, 1)
    var appMgrParms = getAppManagerParams()

    appMgrParms.KubeClient = kubeClient

    appMgr := appmanager.NewManager(&appMgrParms)

    GetNamespaces(appMgr)

    appMgr.K8sVersion = getk8sVersion()

    log.Infof("[INIT] kubernetes version %s", appMgr.K8sVersion)

    intervalFactor := time.Duration(*nodePollInterval)
    log.Infof("pollers NewNodePoller, nodePollInterval: %d, intervalFactor: %d", *nodePollInterval, intervalFactor)
    np := pollers.NewNodePoller(appMgrParms.KubeClient, intervalFactor*time.Second, *nodeLabelSelector)
    err = setupNodePolling(appMgr, np, eventChan, appMgrParms.KubeClient)
    if nil != err {
        log.Fatalf("Required polling utility for node updates failed setup: %v",err)
    }

    //fmt.Printf("%+v\n", appMgr)


}
