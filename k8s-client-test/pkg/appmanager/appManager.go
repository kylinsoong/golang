package appmanager

import (
    "k8s.io/client-go/kubernetes"
    rest "k8s.io/client-go/rest"
    "k8s.io/client-go/util/workqueue"
    "k8s.io/client-go/tools/cache"
)

type Manager struct {
    kubeClient          kubernetes.Interface
    restClientv1        rest.Interface
    vsQueue      workqueue.RateLimitingInterface
    nsQueue      workqueue.RateLimitingInterface
    informersMutex sync.Mutex
    appInformers map[string]*appInformer
    nsInformer cache.SharedIndexInformer
    eventNotifier *EventNotifier
}

type Params struct {
    KubeClient        kubernetes.Interface
    restClient        rest.Interface

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

    return &manager

}
