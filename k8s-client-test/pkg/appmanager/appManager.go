package appmanager

import (
    "sync"

    "k8s.io/client-go/kubernetes"
    rest "k8s.io/client-go/rest"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/util/workqueue"
    "k8s.io/client-go/tools/cache"
)

type Manager struct {
    kubeClient              kubernetes.Interface
    restClientv1            rest.Interface
    restClientv1beta1       rest.Interface
    netClientv1             rest.Interface
    processedResources      map[string]bool
    processedResourcesMutex sync.Mutex
    processedHostPath       ProcessedHostPath
    vsQueue                 workqueue.RateLimitingInterface
    nsQueue                 workqueue.RateLimitingInterface
    informersMutex          sync.Mutex
    appInformers            map[string]*appInformer
    nsInformer              cache.SharedIndexInformer
    eventNotifier           *EventNotifier
    nplStore                map[string]NPLAnnoations
    nplStoreMutex           sync.Mutex
}

type ProcessedHostPath struct {
    sync.Mutex
    processedHostPathMap map[string]metav1.Time
}

type NPLAnnotation struct {
    PodPort  int32  `json:"podPort"`
    NodeIP   string `json:"nodeIP"`
    NodePort int32  `json:"nodePort"`
}

type NPLAnnoations []NPLAnnotation

type appInformer struct {
        namespace        string
        cfgMapInformer   cache.SharedIndexInformer
        svcInformer      cache.SharedIndexInformer
        endptInformer    cache.SharedIndexInformer
        ingInformer      cache.SharedIndexInformer
        routeInformer    cache.SharedIndexInformer
        nodeInformer     cache.SharedIndexInformer
        secretInformer   cache.SharedIndexInformer
        ingClassInformer cache.SharedIndexInformer
        podInformer      cache.SharedIndexInformer
        stopCh           chan struct{}
}

type Params struct {
    KubeClient             kubernetes.Interface
    restClient             rest.Interface
    broadcasterFunc        NewBroadcasterFunc
    ManageConfigMaps       bool
    ManageIngress          bool
}

func NewManager(params *Params) *Manager {

    vsQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "virtual-server-controller")
    nsQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "namespace-controller")

    manager := Manager {
        kubeClient:             params.KubeClient,
        restClientv1:           params.restClient,
	vsQueue:                vsQueue,
        nsQueue:                nsQueue,
        appInformers:           make(map[string]*appInformer),
        eventNotifier:          NewEventNotifier(params.broadcasterFunc),
    }

    manager.processedResources = make(map[string]bool)
    manager.processedHostPath.processedHostPathMap = make(map[string]metav1.Time)
    manager.nplStore = make(map[string]NPLAnnoations)

    return &manager

}
