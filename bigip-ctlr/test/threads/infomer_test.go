package threads

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

var (
	cfgMapInformer cache.SharedIndexInformer
	kubeClient     kubernetes.Interface
	restClientv1   rest.Interface
	vsQueue        workqueue.RateLimitingInterface
)

const (
	Configmaps    = "configmaps"
	OprTypeCreate = "create"
	OprTypeUpdate = "update"
	OprTypeDelete = "delete"
)

type serviceQueueKey struct {
	Namespace    string
	ServiceName  string
	ResourceKind string
	ResourceName string
	Operation    string
}

func checkValidConfigMap(obj interface{}, oprType string) (bool, []*serviceQueueKey) {

	var keyList []*serviceQueueKey
	cm := obj.(*v1.ConfigMap)
	namespace := cm.ObjectMeta.Namespace

	err := validateConfigJson(cm.Data["template"])
	if err != nil {
		fmt.Errorf("Error processing configmap %v in namespace: %v with err: %v", cm.Name, cm.Namespace, err)
		return false, nil
	}

	if ok := processAgentLabels(cm.Labels, cm.Name, namespace); ok {
		key := &serviceQueueKey{
			Namespace:    namespace,
			Operation:    oprType,
			ResourceKind: Configmaps,
			ResourceName: cm.Name,
		}
		keyList = append(keyList, key)
		return true, keyList
	}

	return false, nil
}

func processAgentLabels(m map[string]string, n, ns string) bool {
	funCMapOptions := func(cfg string) bool {
		if cfg == "" {
			return true
		}
		c := strings.Split(cfg, "/")
		if len(c) == 2 {
			if n == c[1] && ns == c[0] {
				return true
			}
			return false
		}
		return true
	}
	if m["overrideAS3"] == "true" || m["overrideAS3"] == "false" {
		return funCMapOptions("")
	} else if m["as3"] == "true" || m["as3"] == "false" {
		return true
	}
	return false
}

func enqueueConfigMap(obj interface{}, operation string) {
	if ok, keys := checkValidConfigMap(obj, operation); ok {
		for _, key := range keys {
			key.Operation = operation
			vsQueue.Add(*key)
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] operation: %s, configmap name: %s, Add %v to vsQueue %v\n", timestamp, key.Operation, key.ResourceName, *key, vsQueue)
		}
	}
}

func Test_Configmap_Informer(t *testing.T) {

	t.Log("Test Start")

	vsQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "virtual-server-controller")

	kubeConfig := "/Users/k.song/src/golang/bigip-ctlr/config"
	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		t.Errorf("Create Kubernetes Config Error %v", err)
	}
	kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("Create Kubernetes Client Error %v", err)
	}

	restClientv1 = kubeClient.CoreV1().RESTClient()

	syncInterval := 30
	resyncPeriod := time.Duration(syncInterval) * time.Second

	namespace := "f5-hub-1"
	everything := func(options *metav1.ListOptions) {
		options.LabelSelector = ""
	}

	cfgMapInformer = cache.NewSharedIndexInformer(
		cache.NewFilteredListWatchFromClient(
			restClientv1,
			Configmaps,
			namespace,
			everything,
		),
		&v1.ConfigMap{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	cfgMapInformer.AddEventHandlerWithResyncPeriod(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) { enqueueConfigMap(obj, OprTypeCreate) },
			UpdateFunc: func(old, cur interface{}) { enqueueConfigMap(cur, OprTypeUpdate) },
			DeleteFunc: func(obj interface{}) { enqueueConfigMap(obj, OprTypeDelete) },
		},
		resyncPeriod,
	)

	go func() {
		defer vsQueue.ShutDown()
		cfgMapInformer.Run(make(chan struct{}))
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	t.Log("Test Exit")

}
