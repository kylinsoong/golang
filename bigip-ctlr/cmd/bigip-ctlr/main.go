package main

import (
    "fmt"
    "strings"
    "os"
    "time"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

    cisAgent "github.com/kylinsoong/bigip-ctlr/pkg/agent"
    log "github.com/kylinsoong/bigip-ctlr/pkg/vlogger"
    clog "github.com/kylinsoong/bigip-ctlr/pkg/vlogger/console"
    "github.com/kylinsoong/bigip-ctlr/pkg/writer"
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
    configWriter       writer.Writer
    agRspChan          chan interface{}
)




func getConfigWriter() writer.Writer {
    if configWriter == nil {
        var err error
        configWriter, err = writer.NewConfigWriter()
        if nil != err {
            log.Fatalf("[INIT] Failed creating ConfigWriter tool: %v", err)
            os.Exit(1)
        }
    }
    return configWriter
} 

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

func initLogger(logLevel string) error {
    log.RegisterLogger(log.LL_MIN_LEVEL, log.LL_MAX_LEVEL, clog.NewConsoleLogger())
                
    if ll := log.NewLogLevel(logLevel); nil != ll {
        log.SetLogLevel(*ll)
    } else {
        return fmt.Errorf("Unknown log level requested: %s\n" + "    Valid log levels are: DEBUG, INFO, WARNING, ERROR, CRITICAL", logLevel)
    }
    return nil 
} 

func main() {

    fmt.Println("F5 BIG-IP Controller Start")

    logLevel := "DEBUG"
    initLogger(logLevel)

    dgPath := strings.Join([]string{resource.DEFAULT_PARTITION, "Shared"}, "/")
    log.Infof("BIG-IP Data Group Path: %s", dgPath)

    appmanager.RegisterBigIPSchemaTypes()

    config, err := getKubeConfig()
    if err != nil {
        os.Exit(1)
    }

    kubeClient, err = kubernetes.NewForConfig(config)
    if err != nil {
        os.Exit(1)
    }

    log.Infof("create kubeClient via %s", config.Host)

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

    log.Infof("Create Global Section %v, BIG-IP Section struct %v", gs, bs)


    pythonBaseDir := ""    
    subPidCh, err := startPythonDriver(getConfigWriter(), gs, bs, pythonBaseDir)
    if nil != err {
        log.Fatalf("Could not initialize subprocess configuration: %v", err)
    }
    log.Infof("Current don't has subpid, %v", subPidCh)
/*
    subPid := <-subPidCh
    defer func(pid int) {
        if 0 != pid {
            var proc *os.Process
            proc, err = os.FindProcess(pid)
            if nil != err {
                log.Warningf("Failed to find sub-process on exit: %v", err)
            }
            err = proc.Signal(os.Interrupt)
            if nil != err {
                log.Warningf("Could not stop sub-process on exit: %d - %v", pid, err)
            }
        }
    }(subPid)
*/

    //agRspChan = make(chan interface{}, 1)
    //var appMgrParms = getAppManagerParams()

    log.Infof("[INIT] Creating Agent for %v", agent)

    nodePollInterval := 30
    intervalFactor := time.Duration(nodePollInterval)

    log.Infof("nodePollInterval: %d, intervalFactor: %d", nodePollInterval, intervalFactor)
    

    log.Infof("Started /metrics and /health service")

    log.Infof("appMgr run")

    

    fmt.Println("F5 BIG-IP Controller End")

}
