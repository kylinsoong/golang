package main

import (
    "fmt"
    "os"
    "strings"
    "os/signal"
    "syscall"
    "net/url"

    "github.com/spf13/pflag"
    "golang.org/x/crypto/ssh/terminal"

    restAgent "github.com/kylinsoong/golang/nim-rest-client/pkg/agent"
    "github.com/kylinsoong/golang/nim-rest-client/pkg/agent/nim"
    log  "github.com/kylinsoong/golang/nim-rest-client/pkg/vlogger"
    clog "github.com/kylinsoong/golang/nim-rest-client/pkg/vlogger/console"
)

var (
    flags           *pflag.FlagSet
    nimFlags        *pflag.FlagSet

    printHelp       *bool
    logLevel        *string
    enableTLS       *string
    ciphers         *string
    restURL         *string
    restUsername    *string
    restPassword    *string
    sslInsecure     *bool

    groups          *string
    instances       *string
    stages          *string    
    
    agent           *string
    AgentREST       restAgent.RESTAgentInterface

)

func _init() {
    flags = pflag.NewFlagSet("main", pflag.PanicOnError)
    nimFlags = pflag.NewFlagSet("Global", pflag.PanicOnError)

    var err error
    var width int
    fd := int(os.Stdout.Fd())
    if terminal.IsTerminal(fd) {
        width, _, err = terminal.GetSize(fd)
        if nil != err {
            width = 0
        }
    }

    printHelp = nimFlags.Bool("help", false, "Optional, print help and exit.")
    logLevel = nimFlags.String("log-level", "INFO", "Optional, logging level")
    enableTLS = nimFlags.String("tls-version", "1.2", "Optional, Configure TLS version to be enabled")
    ciphers = nimFlags.String("ciphers", "DEFAULT", "Optional, Configures a ciphersuite selection string. cipher-group and ciphers are mutually exclusive, only use one.")
    restURL = nimFlags.String("url", "", "Required, URL for the REST")
    restUsername = nimFlags.String("username", "", "Required, user name for the REST user account.")
    restPassword = nimFlags.String("password", "", "Required, password for the REST user account.")
    sslInsecure = nimFlags.Bool("insecure", true, "Optional, when set to true, enable insecure SSL communication to REST.")

    groups = nimFlags.String("instance-group", "", "The name for the instance group")        
    instances = nimFlags.String("instance", "", "The name for the instance")
    stages = nimFlags.String("stage-config", "", "The name for staged config")

    agent = nimFlags.String("agent", "nim", "Optional, allowed value are nim, as3")

    nimFlags.Usage = func() {
        fmt.Fprintf(os.Stderr, "  Flags:\n%s\n", nimFlags.FlagUsagesWrapped(width))
    }

    flags.AddFlagSet(nimFlags)

    flags.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
        nimFlags.Usage()
    }
}

func initLogger(logLevel string) error {
    log.RegisterLogger(log.LL_MIN_LEVEL, log.LL_MAX_LEVEL, clog.NewConsoleLogger())

    if ll := log.NewLogLevel(logLevel); nil != ll {
        log.SetLogLevel(*ll)
    } else {
        return fmt.Errorf("Unknown log level requested: %s\n"+"    Valid log levels are: DEBUG, INFO, WARNING, ERROR, CRITICAL", logLevel)
    }
    return nil
}

func init() {
    _init()
}

func getUserAgentInfo() string {
    return fmt.Sprintf("NIM Client/v%s", "0.1")
}

func getNIMParams() *nim.Params {
    return &nim.Params{
        EnableTLS:               *enableTLS,
        Ciphers:                 *ciphers,
        NIMUsername:             *restUsername,
        NIMPassword:             *restPassword,
	NIMURL:                  *restURL,
        TrustedCerts:            "",
        SSLInsecure:             *sslInsecure,
        UserAgent:               getUserAgentInfo(),
    }
}

func getAgentParams(agent string) interface{} {
    var params interface{}
    switch agent {
    case restAgent.AS3Agent:
        params = getNIMParams()
    case restAgent.NIMAgent:
        params = getNIMParams()
    }
    return params
}

func verifyArgs() error {

    if !strings.HasPrefix(*restURL, "https://") {
        *restURL = "https://" + *restURL
    }
    u, err := url.Parse(*restURL)
    if nil != err {
        return fmt.Errorf("Error parsing url: %s", err)
    }
    if len(u.Path) > 0 && u.Path != "/" {
        return fmt.Errorf("URL path must be empty or '/'; check URL formatting and/or remove %s from path",u.Path)
    }

    return nil
}

func main() {

    flags.Parse(os.Args)

    if *printHelp {
        flags.Usage()
	os.Exit(1)
    } 

    *logLevel = strings.ToUpper(*logLevel)
    initLogger(*logLevel)

    defer func() {
        if r := recover(); r != nil {
            log.Errorf("Recovered panic: %v", r)
        }       
    }() 
 
    err := verifyArgs()
    if nil != err {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	flags.Usage()
	os.Exit(1)
    }

    interruptSignal := make(chan os.Signal, 1)
    signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
	<-interruptSignal
	log.Infof("Received interrupt signal. Shutting down gracefully...")
	os.Exit(0)
    }()

    log.Infof("[INIT] NGINX Instance Manager REST Client")

    AgentREST, err= restAgent.CreateAgent(*agent)
    if err != nil {
        log.Fatalf("[INIT] unable to create agent %v error: err: %+v\n", *agent, err)
        os.Exit(1)
    }
    if err = AgentREST.Init(getAgentParams(*agent)); err != nil {
        log.Fatalf("[INIT] Failed to initialize %v agent, %+v\n", *agent, err)
	os.Exit(1)
    }
    defer AgentREST.DeInit()

    if *groups != "" && len(*groups) > 1 {
        AgentREST.DoBatchDeploy(*groups, *stages)
    }
 
}
