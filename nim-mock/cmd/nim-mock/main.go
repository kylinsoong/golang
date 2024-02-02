package main

import (
    "fmt"
    "os"
    "net/http"
    "strings"
    "crypto/tls"
    "os/signal"
    "syscall"

    "github.com/spf13/pflag"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/gorilla/mux"

    "github.com/kylinsoong/golang/nim-mock/pkg/mock"
    log  "github.com/kylinsoong/golang/nim-mock/pkg/vlogger"
    clog "github.com/kylinsoong/golang/nim-mock/pkg/vlogger/console"
)

var (
    flags           *pflag.FlagSet
    mockFlags       *pflag.FlagSet

    logLevel        *string
    tls_certificate *string
    tls_private     *string
)

func _init() {
    flags = pflag.NewFlagSet("main", pflag.PanicOnError)
    mockFlags = pflag.NewFlagSet("Global", pflag.PanicOnError)

    var err error
    var width int
    fd := int(os.Stdout.Fd())
    if terminal.IsTerminal(fd) {
        width, _, err = terminal.GetSize(fd)
        if nil != err {
            width = 0
        }
    }

    logLevel = mockFlags.String("log-level", "INFO", "Optional, logging level")
    tls_certificate = mockFlags.String("tls-certificate", "/app/certificate.crt", "Optional, tls certificate")
    tls_private = mockFlags.String("tls-private", "/app/private.key", "Optional, tls private key")

    mockFlags.Usage = func() {
        fmt.Fprintf(os.Stderr, "  Flags:\n%s\n", mockFlags.FlagUsagesWrapped(width))
    }

    flags.AddFlagSet(mockFlags)

    flags.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
        mockFlags.Usage()
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

func main() {

    flags.Parse(os.Args)
    *logLevel = strings.ToUpper(*logLevel)
    initLogger(*logLevel)

    defer func() {
        if r := recover(); r != nil {
            log.Errorf("Recovered panic: %v", r)
        }       
    }() 

    interruptSignal := make(chan os.Signal, 1)
    signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
	<-interruptSignal
	log.Infof("Received interrupt signal. Shutting down gracefully...")
	os.Exit(0)
    }()

    r := mux.NewRouter()

    r.HandleFunc("/api/platform/v1/instance-groups/summary", mock.GetInstanceGroupsSummary).Methods("GET")
    r.HandleFunc("/api/platform/v1/instance-groups/{uid}/config", mock.InstanceGroupsConfig).Methods("GET", "POST")

    cert, err := tls.LoadX509KeyPair(*tls_certificate, *tls_private)
    if err != nil {
	log.Errorf("Error loading SSL certificate and private key: %v", err)
    }

    tlsConfig := &tls.Config{
	Certificates: []tls.Certificate{cert},
    }

    server := &http.Server{
	Addr:      ":443",
	TLSConfig: tlsConfig,
        Handler: r,
    }

    log.Infof("Start NIM Mock Service")
 
   err = server.ListenAndServeTLS("", "")
    if err != nil {
	log.Errorf("Error starting server: %v", err)
    }

    log.Infof("AS3 Mock Service Started")
}
