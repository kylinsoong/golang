package threads

import (
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func exampleWork() {
	fmt.Println("Doing some work...")
}

func Test_k8s_apimachinery_pkg_util_wait(t *testing.T) {
	stopCh := make(chan struct{})
	go wait.Until(exampleWork, time.Second, stopCh)
	time.Sleep(5 * time.Second)
	close(stopCh)
	time.Sleep(1 * time.Second)
	fmt.Println("Main goroutine exiting...")
}
