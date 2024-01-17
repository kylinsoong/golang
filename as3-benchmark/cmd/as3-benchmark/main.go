/*
How to run?

go run cmd/as3-benchmark/main.go --ops=del --declaration=$(pwd)/cm1000.txt --bigip-host=192.168.45.52 --bigip-username=admin --bigip-password=admin
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/kylinsoong/as3-benchmark/pkg/as3perf"
	log "github.com/kylinsoong/as3-benchmark/pkg/vlogger"
	clog "github.com/kylinsoong/as3-benchmark/pkg/vlogger/console"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// Flag sets and supported flags
	flags       *pflag.FlagSet
	globalFlags *pflag.FlagSet

	// Custom Resource
	ops           *string
	declaration   *string
	logLevel      *string
	bigIPHost     *string
	bigIPUsername *string
	bigIPPassword *string
)

func _init() {
	flags = pflag.NewFlagSet("main", pflag.PanicOnError)
	globalFlags = pflag.NewFlagSet("Global", pflag.PanicOnError)

	// Flag wrapping
	var err error
	var width int
	fd := int(os.Stdout.Fd())
	if terminal.IsTerminal(fd) {
		width, _, err = terminal.GetSize(fd)
		if nil != err {
			width = 0
		}
	}

	logLevel = globalFlags.String("log-level", "INFO", "Optional, logging level")
	ops = globalFlags.String("ops", "add", "Set the operation, default is add, supported: add, del")
	declaration = globalFlags.String("declaration", "", "Set the as3 declaration file")
	bigIPHost = globalFlags.String("bigip-host", "", "Required, Host for the Big-IP")
	bigIPUsername = globalFlags.String("bigip-username", "", "Required, user name for the Big-IP user account.")
	bigIPPassword = globalFlags.String("bigip-password", "", "Required, password for the Big-IP user account.")

	globalFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "  Flags:\n%s\n", globalFlags.FlagUsagesWrapped(width))
	}

	flags.AddFlagSet(globalFlags)

	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
		globalFlags.Usage()
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

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	err := flags.Parse(os.Args)
	if nil != err {
		os.Exit(1)
	}

	*logLevel = strings.ToUpper(*logLevel)
	initLogger(*logLevel)

	log.Infof("[MAIN] Starting: AS3 Performance Benchmark")
	log.Debugf("ops: %s, declaration: %s, bigIPHost: %s, bigIPUsername: %s, bigIPPassword: %s", *ops, *declaration, *bigIPHost, *bigIPUsername, *bigIPPassword)

	startTime := time.Now()
	if *ops == "add" {
		Addition(*declaration, *bigIPHost, *bigIPUsername, *bigIPPassword)
	} else {
		Deletion(*declaration, *bigIPHost, *bigIPUsername, *bigIPPassword)
	}
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	log.Infof("[MAIN] Total time taken: %s", totalTime)

}
