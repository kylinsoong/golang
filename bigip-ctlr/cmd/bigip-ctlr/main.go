package main

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "os"
    "time"
    "os/signal"
    "syscall"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/labels"

    "github.com/spf13/pflag"
    "golang.org/x/crypto/ssh/terminal"

    "github.com/kylinsoong/bigip-ctlr/pkg/appmanager"
    "github.com/kylinsoong/bigip-ctlr/pkg/pollers"
    "github.com/kylinsoong/bigip-ctlr/pkg/resource"
    log "github.com/kylinsoong/bigip-ctlr/pkg/vlogger"
    clog "github.com/kylinsoong/bigip-ctlr/pkg/vlogger/console"
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
    syncInterval     *int

    namespaceLabel         *string
    namespaces             *[]string
    useNodeInternal        *bool
    poolMemberType         *string
    inCluster              *bool
    kubeConfig             *string
    manageConfigMaps       *bool
    manageIngress          *bool
    hubMode                *bool
    nodeLabelSelector      *string
    watchAllNamespaces     bool
    isNodePort             bool

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
    syncInterval = globalFlags.Int("periodic-sync-interval", 30, "Optional, interval (in seconds) at which to queue resources.")

    globalFlags.Usage = func() {
        fmt.Fprintf(os.Stderr, "  Global:\n%s\n", globalFlags.FlagUsagesWrapped(width))
    }

    namespaceLabel = kubeFlags.String("namespace-label", "", "Optional, used to watch for namespaces with this label")
    namespaces = kubeFlags.StringArray("namespace", []string{}, "Optional, Kubernetes namespace(s) to watch. If left blank controller will watch all k8s namespaces")
    useNodeInternal = kubeFlags.Bool("use-node-internal", true, "Optional, provide kubernetes InternalIP addresses to pool")
    poolMemberType = kubeFlags.String("pool-member-type", "nodeport",
                "Optional, type of BIG-IP pool members to create. "+
                        "'nodeport' will use k8s service NodePort. "+
                        "'cluster' will use service endpoints. "+
                        "The BIG-IP must be able access the cluster network"+
                        "'nodeportlocal' only supported with antrea cni")
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
        UseNodeInternal:        *useNodeInternal,
        IsNodePort:             isNodePort,
        ManageConfigMaps:       *manageConfigMaps,
        ManageIngress:          *manageIngress,
        AgRspChan:              agRspChan,
        HubMode:                *hubMode,
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

    err := np.RegisterListener(appMgr.ProcessNodeUpdate)
    if nil != err {
        return fmt.Errorf("error registering node update listener: %v", err)
    }

    return nil
}

func createLabel(label string) (labels.Selector, error) {
    var l labels.Selector
    var err error
    if label == "" {
       l = labels.Everything()
    } else {
        l, err = labels.Parse(label)
        if err != nil {
            return nil, fmt.Errorf("failed to parse Label Selector string: %v", err)
        }
    }
    return l, nil
}

func setupWatchers(appMgr *appmanager.Manager, resyncPeriod time.Duration)  {
    label := resource.DefaultConfigMapLabel

    if len(*namespaceLabel) == 0  {
        ls, err := createLabel("")
        if nil != err {
            log.Warningf("[INIT] Failed to create label selector: %v", err)
        }

        err = appMgr.AddNamespaceLabelInformer(ls, resyncPeriod)
        if nil != err {
            log.Warningf("[INIT] Failed to add label watch for all namespaces:%v", err)
        }

        ls, err = createLabel(label)
        if nil != err {
            log.Warningf("[INIT] Failed to create label selector: %v", err)
        }

        if watchAllNamespaces == true {
            err = appMgr.AddNamespace("", ls, resyncPeriod)
            if nil != err {
                log.Warningf("[INIT] Failed to add informers for all namespaces:%v", err)
            }
        } else {
            for _, namespace := range *namespaces {
                err = appMgr.AddNamespace(namespace, ls, resyncPeriod)
                if nil != err {
                    log.Warningf("[INIT] Failed to add informers for namespace %v: %v", namespace, err)
                } else {
                    log.Debugf("[INIT] Added informers for namespace %v", namespace)
                }
            }
        }

    } else {
        ls, err := createLabel(*namespaceLabel)
        if nil != err {
            log.Warningf("[INIT] Failed to create label selector: %v", err)
        }
        err = appMgr.AddNamespaceLabelInformer(ls, resyncPeriod)
        if nil != err {
            log.Warningf("[INIT] Failed to add label watch for all namespaces:%v", err)
        }
        appMgr.DynamicNS = true
    }
}

func main() {

    err := flags.Parse(os.Args)
    if nil != err {
        os.Exit(1)
    }

    *logLevel = strings.ToUpper(*logLevel)
    initLogger(*logLevel)

    if len(*namespaces) == 0 && len(*namespaceLabel) == 0 {
        watchAllNamespaces = true
    } else {
        watchAllNamespaces = false
    }

    if *poolMemberType == "nodeport" {
        isNodePort = true
    } else if *poolMemberType == "cluster" || *poolMemberType == "nodeportlocal" {
        isNodePort = false
    } else {
        return fmt.Errorf("'%v' is not a valid Pool Member Type", *poolMemberType)
    }

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
    log.Infof("Setup node poller, nodePollInterval: %d", *nodePollInterval)
    np := pollers.NewNodePoller(appMgrParms.KubeClient, intervalFactor*time.Second, *nodeLabelSelector)
    err = setupNodePolling(appMgr, np, eventChan, appMgrParms.KubeClient)
    if nil != err {
        log.Fatalf("Required polling utility for node updates failed setup: %v",err)
    }

    np.Run()
    defer np.Stop()

    log.Infof("Setup watchers, syncInterval: %d", *syncInterval)
    setupWatchers(appMgr, time.Duration(*syncInterval)*time.Second)

    stopCh := make(chan struct{})

    log.Infof("appMgr run")
    appMgr.Run(stopCh)

    log.Infof("signal process")
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    sig := <-sigs
    close(stopCh) 
    log.Infof("[INIT] Exiting - signal %v\n", sig)
    //fmt.Printf("%+v\n", appMgr)


}
