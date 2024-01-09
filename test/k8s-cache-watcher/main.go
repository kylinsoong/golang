package main


import (
    "fmt"
    "time"
    "os"
    "runtime"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"    

    "k8s.io/client-go/util/workqueue"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"

    v1 "k8s.io/api/core/v1"
)

const (
    Namespaces     = "namespaces"
    Services       = "services"
    Endpoints      = "endpoints"
    Configmaps     = "configmaps"
    Ingresses      = "ingresses"
    syncInterval   = 30 * time.Second 

    OprTypeCreate   = "create"
    OprTypeUpdate  = "update"
    OprTypeDelete  = "delete"
)

var (
    restClientv1     rest.Interface
    kubeClient       kubernetes.Interface

    vsQueue          workqueue.RateLimitingInterface
    nsQueue          workqueue.RateLimitingInterface

    nsInformer       cache.SharedIndexInformer
    cfgMapInformer   cache.SharedIndexInformer
    svcInformer      cache.SharedIndexInformer
    endptInformer    cache.SharedIndexInformer
    ingInformer      cache.SharedIndexInformer
)

func printerr(err error) {
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1) // exit with a non-zero status to indicate an error
    }
}

func main() {

    vsQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "virtual-server-controller")
    nsQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "namespace-controller")

    var err error

    syncInterval := 30
    resyncPeriod := time.Duration(syncInterval)*time.Second

    label := "cis_scanner=zone1"
    var ls labels.Selector
    ls, err = labels.Parse(label)    
    printerr(err)

    kubeConfig := "/Users/k.song/src/golang/bigip-ctlr/config"
    var config *rest.Config
    config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
    printerr(err)
    kubeClient, err = kubernetes.NewForConfig(config)
    printerr(err)

    restClientv1 = kubeClient.CoreV1().RESTClient()

    optionsModifier := func(options *metav1.ListOptions) {
        options.LabelSelector = ls.String()
    }

    nsInformer = cache.NewSharedIndexInformer(
        cache.NewFilteredListWatchFromClient(
            restClientv1,
            Namespaces,
            "",
            optionsModifier,
        ),
        &v1.Namespace{},
        resyncPeriod,
        cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
    )

    nsInformer.AddEventHandlerWithResyncPeriod(
        &cache.ResourceEventHandlerFuncs{
            AddFunc:    func(obj interface{}) { enqueueNamespace(obj) },
            UpdateFunc: func(old, cur interface{}) { enqueueNamespace(cur) },
            DeleteFunc: func(obj interface{}) { enqueueNamespace(obj) },
        },
        resyncPeriod,
    )

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
	    AddFunc: func(obj interface{}) { enqueueConfigMap(obj, OprTypeCreate) },
	    UpdateFunc: func(old, cur interface{}) {enqueueConfigMap(cur, OprTypeUpdate)},
	    DeleteFunc: func(obj interface{}) { enqueueConfigMap(obj, OprTypeDelete) },
	},
	resyncPeriod,
    )

    svcInformer = cache.NewSharedIndexInformer(
        cache.NewFilteredListWatchFromClient(
            restClientv1,
            Services,
            namespace,
            everything,
        ),
        &v1.Service{},
        resyncPeriod,
        cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
    )

    svcInformer.AddEventHandlerWithResyncPeriod(
	&cache.ResourceEventHandlerFuncs{
	    AddFunc:    func(obj interface{}) { enqueueService(obj, OprTypeCreate) },
	    UpdateFunc: func(old, cur interface{}) { enqueueService(cur, OprTypeUpdate) },
	    DeleteFunc: func(obj interface{}) { enqueueService(obj, OprTypeDelete) },
	},
	resyncPeriod,
    )

    endptInformer = cache.NewSharedIndexInformer(
        cache.NewFilteredListWatchFromClient(
            restClientv1,
            Endpoints,
            namespace,
            everything,
        ),
        &v1.Endpoints{},
        resyncPeriod,
        cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
    )

    endptInformer.AddEventHandlerWithResyncPeriod(
	&cache.ResourceEventHandlerFuncs{
	    AddFunc:    func(obj interface{}) { enqueueEndpoints(obj, OprTypeCreate) },
	    UpdateFunc: func(old, cur interface{}) { enqueueEndpoints(cur, OprTypeUpdate) },
	    DeleteFunc: func(obj interface{}) { enqueueEndpoints(obj, OprTypeDelete) },
	},
	resyncPeriod,
    )

    go func() {
        defer nsQueue.ShutDown()
	nsInformer.Run(make(chan struct{}))
    }()

    go func() {
        defer vsQueue.ShutDown()
        cfgMapInformer.Run(make(chan struct{}))
        svcInformer.Run(make(chan struct{}))
        endptInformer.Run(make(chan struct{}))
    }()

    go func() {
        for {
	    processNextNamespace()
	}
    }()

    go func() {
        for {
            processNextVS()
        }
    }()

    select {}
}

func enqueueNamespace(obj interface{}) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    ns := obj.(*v1.Namespace)
    nsQueue.Add(ns.ObjectMeta.Name)
    fmt.Printf("[%s] (%d) add %s to queue\n", timestamp, getGoroutineID(), ns.ObjectMeta.Name)
}

func enqueueConfigMap(obj interface{}, operation string) {
    cm := obj.(*v1.ConfigMap)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    vsQueue.Add(cm.ObjectMeta.Name)
    fmt.Printf("[%s] (%d) %s configmap %s, add to queue\n", timestamp, getGoroutineID(), operation, cm.ObjectMeta.Name)
}

func enqueueService(obj interface{}, operation string) {
    svc := obj.(*v1.Service)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    vsQueue.Add(svc.ObjectMeta.Name)
    fmt.Printf("[%s] (%d) %s service %s, add to queue\n", timestamp, getGoroutineID(), operation, svc.ObjectMeta.Name)
}

func  enqueueEndpoints(obj interface{}, operation string) {
    eps := obj.(*v1.Endpoints)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    vsQueue.Add(eps.ObjectMeta.Name)
    fmt.Printf("[%s] (%d) %s endpoint %s, add to queue\n", timestamp, getGoroutineID(), operation, eps.ObjectMeta.Name)
}

func processNextNamespace() {
    item, quit := nsQueue.Get()
    if quit {
        return
    }
    defer nsQueue.Done(item)

    namespaceName, ok := item.(string)
    if !ok {
        fmt.Println("Failed to process namespace, unexpected item type")
	return
    }
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] (%d) processing Namespace: %s\n", timestamp, getGoroutineID(), namespaceName)
}

func processNextVS() {
    item, quit := vsQueue.Get()
    if quit {
        return
    }
    defer vsQueue.Done(item)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] (%d) processing %s\n", timestamp, getGoroutineID(), item)
}

func getGoroutineID() int64 {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getGoroutineID:", r)
		}
	}()

	buf := make([]byte, 64)
	runtime.Stack(buf, false)
	var id int64
	n, _ := fmt.Sscanf(string(buf), "goroutine %d", &id)
	if n != 1 {
		fmt.Println("Failed to extract Goroutine ID")
		return -1
	}
	return id
}
