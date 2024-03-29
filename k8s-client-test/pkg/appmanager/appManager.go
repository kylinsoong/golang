package appmanager

import (
    "sync"
    "fmt"
    "reflect"
    "time"
    "context"
    "strings"

    "k8s.io/client-go/kubernetes"

    netv1 "k8s.io/api/networking/v1"
    rest "k8s.io/client-go/rest"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    v1 "k8s.io/api/core/v1"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"

    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/client-go/util/workqueue"
    "k8s.io/client-go/tools/cache"

    cisAgent "github.com/kylinsoong/k8s-client-test/pkg/agent"
    log "github.com/kylinsoong/k8s-client-test/pkg/vlogger"
    . "github.com/kylinsoong/k8s-client-test/pkg/resource"
)

type ResourceMap map[int32][]*ResourceConfig

type Manager struct {
    resources               *Resources
    agentCfgMap             map[string]*AgentCfgMap
    kubeClient              kubernetes.Interface
    restClientv1            rest.Interface
    restClientv1beta1       rest.Interface
    netClientv1             rest.Interface
    steadyState             bool
    queueLen                int
    processedItems          int
    processedResources      map[string]bool
    processedResourcesMutex sync.Mutex
    processedHostPath       ProcessedHostPath
    useNodeInternal         bool
    isNodePort              bool
    hubMode                 bool
    oldNodesMutex           sync.Mutex
    oldNodes                []Node
    vsQueue                 workqueue.RateLimitingInterface
    nsQueue                 workqueue.RateLimitingInterface
    informersMutex          sync.Mutex
    appInformers            map[string]*appInformer
    as3Informer             *appInformer
    nsInformer              cache.SharedIndexInformer
    manageConfigMaps        bool
    manageIngress           bool
    DynamicNS               bool
    eventNotifier           *EventNotifier
    nplStore                map[string]NPLAnnoations
    nplStoreMutex           sync.Mutex
    AgentCIS                cisAgent.CISAgentInterface
    agRspChan               chan interface{}
    WatchedNS               WatchedNamespaces
    configMapLabel          string
    K8sVersion              string
    rsrcSSLCtxt             map[string]*v1.Secret
    mergedRulesMap          map[string]map[string]MergedRuleEntry
}

type WatchedNamespaces struct {
    Namespaces     []string
    NamespaceLabel string
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

type serviceQueueKey struct {
    Namespace    string
    ServiceName  string
    ResourceKind string
    ResourceName string
    Operation    string
}

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
    UseNodeInternal        bool
    IsNodePort             bool
    broadcasterFunc        NewBroadcasterFunc
    steadyState            bool
    ManageConfigMaps       bool
    ManageIngress          bool
    HubMode                bool
    AgRspChan              chan interface{}
}

type Node struct {
    Name string
    Addr string
}

func (s *vsSyncStats) isStatsAvailable() bool {
    switch {
    case s.vsUpdated > 0,
         s.vsDeleted > 0,
         s.cpUpdated > 0,
         s.dgUpdated > 0,
         s.poolsUpdated > 0:
         return true
    }

    return false
}

type vsSyncStats struct {
    vsFound      int
    vsUpdated    int
    vsDeleted    int
    cpUpdated    int
    dgUpdated    int
    poolsUpdated int
}

const (
    Namespaces     = "namespaces"
    Services       = "services"
    Endpoints      = "endpoints"
    Secrets        = "secrets"
    Configmaps     = "configmaps"
    Ingresses      = "ingresses"
    IngressClasses = "ingressclasses"
    hubModeInterval  = 30 * time.Second //Hubmode ConfigMap resync interval
)

func NewManager(params *Params) *Manager {

    vsQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "virtual-server-controller")
    nsQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "namespace-controller")

    manager := Manager {
        resources:              NewResources(),
        kubeClient:             params.KubeClient,
        restClientv1:           params.restClient,
        useNodeInternal:        params.UseNodeInternal,
        steadyState:            params.steadyState,
        isNodePort:             params.IsNodePort,
        queueLen:               0,
        processedItems:         0,
	vsQueue:                vsQueue,
        nsQueue:                nsQueue,
        appInformers:           make(map[string]*appInformer),
        eventNotifier:          NewEventNotifier(params.broadcasterFunc),
        hubMode:                params.HubMode,
        manageIngress:          params.ManageIngress,
        manageConfigMaps:       params.ManageConfigMaps,
        agRspChan:              params.AgRspChan,
        agentCfgMap:            make(map[string]*AgentCfgMap),
        rsrcSSLCtxt:            make(map[string]*v1.Secret),
        mergedRulesMap:         make(map[string]map[string]MergedRuleEntry),
    }

    manager.processedResources = make(map[string]bool)
    manager.processedHostPath.processedHostPathMap = make(map[string]metav1.Time)
    manager.nplStore = make(map[string]NPLAnnoations)

    go manager.agentResponseWorker()

    if nil != manager.kubeClient && nil == manager.restClientv1 {
        manager.restClientv1 = manager.kubeClient.CoreV1().RESTClient()
    }

    if nil != manager.kubeClient && nil == manager.restClientv1beta1 {
        manager.restClientv1beta1 = manager.kubeClient.ExtensionsV1beta1().RESTClient()
    }

    if nil != manager.kubeClient && nil == manager.netClientv1 {
        manager.netClientv1 = manager.kubeClient.NetworkingV1().RESTClient()
    }

    return &manager

}

func (appMgr *Manager) getNodes(obj interface{},) ([]Node, error) {
    nodes, ok := obj.([]v1.Node)
    if false == ok {
        return nil, fmt.Errorf("poll update unexpected type, interface is not []v1.Node")
    }

    watchedNodes := []Node{}

    var addrType v1.NodeAddressType
    if appMgr.UseNodeInternal() {
        addrType = v1.NodeInternalIP
    } else {
        addrType = v1.NodeExternalIP
    }

        // Append list of nodes to watchedNodes
    for _, node := range nodes {
        nodeAddrs := node.Status.Addresses
        for _, addr := range nodeAddrs {
            if addr.Type == addrType {
                n := Node{
                    Name: node.ObjectMeta.Name,
                    Addr: addr.Address,
                }
                watchedNodes = append(watchedNodes, n)
            }
        }
    }

    return watchedNodes, nil
}

func (appMgr *Manager) ProcessNodeUpdate( obj interface{}, err error,) {

    if nil != err {
        log.Warningf("[CORE] Unable to get list of nodes, err=%+v", err)
        return
    }

    newNodes, err := appMgr.getNodes(obj)
    if nil != err {
        log.Warningf("[CORE] Unable to get list of nodes, err=%+v", err)
        return
    }

    if appMgr.steadyState {
        if !reflect.DeepEqual(newNodes, appMgr.oldNodes) {
            log.Infof("[CORE] ProcessNodeUpdate: Change in Node state detected")
            items := make(map[serviceQueueKey]int)
            appMgr.resources.ForEach(func(key ServiceKey, cfg *ResourceConfig) {
                queueKey := serviceQueueKey{
                    Namespace:   key.Namespace,
                    ServiceName: key.ServiceName,
                }
                items[queueKey]++
            })

            for queueKey := range items {
                appMgr.vsQueue.Add(queueKey)
            }

            appMgr.oldNodes = newNodes
            log.Warningf("[CORE] Update node cache, %v", appMgr.oldNodes)
        }
    } else {
        appMgr.oldNodes = newNodes
        log.Warningf("[CORE] Initialize appMgr nodes on our first pass through, %v", appMgr.oldNodes)
    }

}

func (appMgr *Manager) Run(stopCh <-chan struct{}) {
    go appMgr.runImpl(stopCh)
}

func (appMgr *Manager) runImpl(stopCh <-chan struct{}) {
    defer utilruntime.HandleCrash()
    defer appMgr.vsQueue.ShutDown()
    defer appMgr.nsQueue.ShutDown()

    if nil != appMgr.nsInformer {
                // Using one worker for namespace label changes.
        appMgr.startAndSyncNamespaceInformer(stopCh)
        go wait.Until(appMgr.namespaceWorker, time.Second, stopCh)
    }

    appMgr.startAndSyncAppInformers()

        // Using only one virtual server worker currently.
    go wait.Until(appMgr.virtualServerWorker, time.Second, stopCh)

    <-stopCh
    appMgr.stopAppInformers()
}

func (appMgr *Manager) startAndSyncNamespaceInformer(stopCh <-chan struct{}) {
    appMgr.informersMutex.Lock()
    defer appMgr.informersMutex.Unlock()
    go appMgr.nsInformer.Run(stopCh)
    cache.WaitForCacheSync(stopCh, appMgr.nsInformer.HasSynced)
}

func (appMgr *Manager) startAndSyncAppInformers() {
    appMgr.informersMutex.Lock()
    defer appMgr.informersMutex.Unlock()
    appMgr.startAppInformersLocked()
    appMgr.waitForCacheSyncLocked()
}

func (appMgr *Manager) startAppInformersLocked() {
    for _, appInf := range appMgr.appInformers {
        appInf.start()
    }
    if nil != appMgr.as3Informer {
        appMgr.as3Informer.start()
    }
}

func (appMgr *Manager) namespaceWorker() {
    for appMgr.processNextNamespace() {
    }
}

func (appMgr *Manager) virtualServerWorker() {
    for appMgr.processNextVirtualServer() {
    }
}

func (appMgr *Manager) stopAppInformers() {
    appMgr.informersMutex.Lock()
    defer appMgr.informersMutex.Unlock()
    for _, appInf := range appMgr.appInformers {
        appInf.stopInformers()
    }
    if nil != appMgr.as3Informer {
        appMgr.as3Informer.stopInformers()
    }
}

func (appMgr *Manager) waitForCacheSyncLocked() {
    for _, appInf := range appMgr.appInformers {
        appInf.waitForCacheSync()
    }
    if nil != appMgr.as3Informer {
        appMgr.as3Informer.waitForCacheSync()
    }
}

func (appInf *appInformer) waitForCacheSync() {
    cacheSyncs := []cache.InformerSynced{}

    if nil != appInf.svcInformer {
        cacheSyncs = append(cacheSyncs, appInf.svcInformer.HasSynced)
    }
    if nil != appInf.endptInformer {
        cacheSyncs = append(cacheSyncs, appInf.endptInformer.HasSynced)
    }
    if nil != appInf.ingInformer {
        cacheSyncs = append(cacheSyncs, appInf.ingInformer.HasSynced)
    }
    if nil != appInf.cfgMapInformer {
        cacheSyncs = append(cacheSyncs, appInf.cfgMapInformer.HasSynced)
    }
    if nil != appInf.nodeInformer {
        cacheSyncs = append(cacheSyncs, appInf.nodeInformer.HasSynced)
    }
    if nil != appInf.ingClassInformer {
        cacheSyncs = append(cacheSyncs, appInf.ingClassInformer.HasSynced)
    }
    cache.WaitForCacheSync(
        appInf.stopCh,
        cacheSyncs...,
    )
}

func (appInf *appInformer) stopInformers() {
    close(appInf.stopCh)
}

func (appInf *appInformer) start() {
    if nil != appInf.svcInformer {
        go appInf.svcInformer.Run(appInf.stopCh)
    }
    if nil != appInf.endptInformer {
        go appInf.endptInformer.Run(appInf.stopCh)
    }
    if nil != appInf.ingInformer {
        go appInf.ingInformer.Run(appInf.stopCh)
    }
    if nil != appInf.routeInformer {
        go appInf.routeInformer.Run(appInf.stopCh)
    }
    if nil != appInf.cfgMapInformer {
        go appInf.cfgMapInformer.Run(appInf.stopCh)
    }
    if nil != appInf.nodeInformer {
        go appInf.nodeInformer.Run(appInf.stopCh)
    }
    if nil != appInf.ingClassInformer {
        go appInf.ingClassInformer.Run(appInf.stopCh)
    }
}

func (appMgr *Manager) processNextNamespace() bool {
    key, quit := appMgr.nsQueue.Get()
    if quit {
        return false
    }
    defer appMgr.nsQueue.Done(key)

    err := appMgr.syncNamespace(key.(string))
    if err == nil {
        appMgr.nsQueue.Forget(key)
        return true
    }

    utilruntime.HandleError(fmt.Errorf("Sync %v failed with %v", key, err))
    appMgr.nsQueue.AddRateLimited(key)

    return true
}

func (appMgr *Manager) processNextVirtualServer() bool {
    key, quit := appMgr.vsQueue.Get()
    if !appMgr.steadyState && appMgr.processedItems == 0 {
        appMgr.queueLen = appMgr.getQueueLength()
    }
    if quit {
                // The controller is shutting down.
        return false
    }

    defer appMgr.vsQueue.Done(key)
    skey := key.(serviceQueueKey)
    if !appMgr.steadyState && !isNonPerfResource(skey.ResourceKind) {
        if skey.Operation != OprTypeCreate {
            appMgr.vsQueue.AddRateLimited(key)
        }
        appMgr.vsQueue.Forget(key)
        return true
    }

    if !appMgr.steadyState && skey.Operation != OprTypeCreate {
        appMgr.vsQueue.AddRateLimited(key)
        appMgr.vsQueue.Forget(key)
        return true
    }

    err := appMgr.syncVirtualServer(skey)
    if err == nil {
        if !appMgr.steadyState {
            appMgr.processedItems++
        }
        appMgr.vsQueue.Forget(key)
        return true
    }

    utilruntime.HandleError(fmt.Errorf("Sync %v failed with %v", key, err))
    appMgr.vsQueue.AddRateLimited(key)

    return true
}

func isNonPerfResource(resKind string) bool {

    switch resKind {
    case Services, Configmaps:
                // Configmaps and Routes get processed according to low performing algorithm
                // But, Service must be processed everytime
        return true
    case Ingresses, Endpoints:
                // Ingresses get processed according to new high performance algorithm
                // Endpoints are out of equation, during initial state never gets processed
        return false
    }

        // Unknown resources are to be considered as non-performing
    return true
}

func (appMgr *Manager) getQueueLength() int {
    qLen := 0

    cmOptions := metav1.ListOptions{
        LabelSelector: appMgr.configMapLabel,
    }

    for _, ns := range appMgr.GetAllWatchedNamespaces() {
        services, err := appMgr.kubeClient.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
        for _, svc := range services.Items {
            if ok, _ := appMgr.checkValidService(&svc); ok {
                qLen++
            }
        }
        if err != nil {
            log.Errorf("[CORE] Failed getting Services from watched namespace : %v.", err)
            return appMgr.vsQueue.Len()
        }

        if false != appMgr.manageConfigMaps {
            cms, err := appMgr.kubeClient.CoreV1().ConfigMaps(ns).List(context.TODO(), cmOptions)
            for _, cm := range cms.Items {
                if ok, _ := appMgr.checkValidConfigMap(&cm, OprTypeCreate); ok {
                    qLen++
                }
            }
            if err != nil {
                log.Errorf("[CORE] Failed getting Configmaps from watched namespace : %v.", err)
                return appMgr.vsQueue.Len()
            }
        }

    }

    return qLen
}

func (appMgr *Manager) GetAllWatchedNamespaces() []string {
    var namespaces []string
    switch {
    case len(appMgr.WatchedNS.Namespaces) != 0:
        namespaces = appMgr.WatchedNS.Namespaces
    case len(appMgr.WatchedNS.NamespaceLabel) != 0:
        NsListOptions := metav1.ListOptions{
            LabelSelector: appMgr.WatchedNS.NamespaceLabel,
        }
        nsL, err := appMgr.kubeClient.CoreV1().Namespaces().List(context.TODO(), NsListOptions)
        if err != nil {
            log.Errorf("[CORE] Error getting Namespaces with Namespace label - %v.", err)
        }
        for _, v := range nsL.Items {
            namespaces = append(namespaces, v.Name)
        }
    default:
        namespaces = append(namespaces, "")
    }
    return namespaces
}


func (appMgr *Manager) UseNodeInternal() bool {
    return appMgr.useNodeInternal
}

func (appMgr *Manager) IsNodePort() bool {
        return appMgr.isNodePort        
}  

func (appMgr *Manager) AddNamespaceLabelInformer(labelSelector labels.Selector, resyncPeriod time.Duration,) error {

    appMgr.informersMutex.Lock()
    defer appMgr.informersMutex.Unlock()

    if nil != appMgr.nsInformer {
        return fmt.Errorf("Already have a namespace label informer added.")
    }

    if 0 != len(appMgr.appInformers) {
        return fmt.Errorf("Cannot set a namespace label informer when informers have been setup for one or more namespaces.")
    }

    optionsModifier := func(options *metav1.ListOptions) {
        options.LabelSelector = labelSelector.String()
    }

    appMgr.nsInformer = cache.NewSharedIndexInformer(
        cache.NewFilteredListWatchFromClient(
            appMgr.restClientv1,
            Namespaces,
            "",
            optionsModifier,
        ),
        &v1.Namespace{},
        resyncPeriod,
        cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
    )

    appMgr.nsInformer.AddEventHandlerWithResyncPeriod(
        &cache.ResourceEventHandlerFuncs{
            AddFunc:    func(obj interface{}) { appMgr.enqueueNamespace(obj) },
            UpdateFunc: func(old, cur interface{}) { appMgr.enqueueNamespace(cur) },
            DeleteFunc: func(obj interface{}) { appMgr.enqueueNamespace(obj) },
        },
        resyncPeriod,
    )

    return nil
}

func (appMgr *Manager) enqueueNamespace(obj interface{}) {
    ns := obj.(*v1.Namespace)
    if !appMgr.DynamicNS && !appMgr.watchingAllNamespacesLocked() {
        if _, ok := appMgr.getNamespaceInformer(ns.Name); !ok {
            return
        }
    }

    appMgr.nsQueue.Add(ns.ObjectMeta.Name)
}

func (appMgr *Manager) watchingAllNamespacesLocked() bool {
    if 0 == len(appMgr.appInformers) {
        // Not watching any namespaces.
        return false
    }
    _, watchingAll := appMgr.appInformers[""]
    return watchingAll
}

func (appMgr *Manager) getNamespaceInformer(ns string,) (*appInformer, bool) {
    appMgr.informersMutex.Lock()
    defer appMgr.informersMutex.Unlock()
    appInf, found := appMgr.getNamespaceInformerLocked(ns)
    return appInf, found
}

func (appMgr *Manager) getNamespaceInformerLocked(ns string,) (*appInformer, bool) {
    toFind := ns
    if appMgr.watchingAllNamespacesLocked() {
        toFind = ""
    }
    appInf, found := appMgr.appInformers[toFind]
    return appInf, found
}

func (appMgr *Manager) AddNamespace(namespace string, cfgMapSelector labels.Selector, resyncPeriod time.Duration,) error {
    appMgr.informersMutex.Lock()
    defer appMgr.informersMutex.Unlock()
    _, err := appMgr.addNamespaceLocked(namespace, cfgMapSelector, resyncPeriod)
    return err
}

func (appMgr *Manager) addNamespaceLocked(namespace string, cfgMapSelector labels.Selector, resyncPeriod time.Duration,) (*appInformer, error) {
        // Check if watching all namespaces by checking all appInformers is created for "" namespace
    if appMgr.watchingAllNamespacesLocked() {
        return nil, fmt.Errorf("Cannot add additional namespaces when already watching all.")
    }

    if len(appMgr.appInformers) > 0 && "" == namespace {
        return nil, fmt.Errorf("Cannot watch all namespaces when already watching specific ones.")
    }

    var appInf *appInformer
    var found bool
    if appInf, found = appMgr.appInformers[namespace]; found {
        return appInf, nil
    }
    appInf = appMgr.newAppInformer(namespace, cfgMapSelector, resyncPeriod)
    appMgr.appInformers[namespace] = appInf
    return appInf, nil
}

func (appMgr *Manager) newAppInformer(namespace string, cfgMapSelector labels.Selector, resyncPeriod time.Duration,) *appInformer {
    log.Debugf("[CORE] Creating new app informer, namespace: %s, cfgMapSelector: %s, resyncPeriod: %d", namespace, cfgMapSelector, resyncPeriod)
    everything := func(options *metav1.ListOptions) {
        options.LabelSelector = ""
    }
    appInf := appInformer{
        namespace: namespace,
        stopCh:    make(chan struct{}),
        svcInformer: cache.NewSharedIndexInformer(
            cache.NewFilteredListWatchFromClient(
                appMgr.restClientv1,
                Services,
                namespace,
                everything,
            ),
            &v1.Service{},
            resyncPeriod,
            cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
        ),
        endptInformer: cache.NewSharedIndexInformer(
            cache.NewFilteredListWatchFromClient(
                appMgr.restClientv1,
                Endpoints,
                namespace,
                everything,
            ),
            &v1.Endpoints{},
            resyncPeriod,
            cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
        ),
    }

    if true == appMgr.manageIngress {
        log.Infof("[CORE] Watching Ingress resources.")
        appInf.ingInformer = cache.NewSharedIndexInformer(
            cache.NewFilteredListWatchFromClient(
                appMgr.netClientv1,
                Ingresses,
                namespace,
                everything,
            ),
            &netv1.Ingress{},
            resyncPeriod,
            cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
        )
        appInf.ingClassInformer = cache.NewSharedIndexInformer(
            cache.NewFilteredListWatchFromClient(
                appMgr.netClientv1,
                IngressClasses,
                "",
                everything,
            ),
            &netv1.IngressClass{},
            resyncPeriod,
            cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
        )
    } else {
        log.Infof("[CORE] Not watching Ingress resources.")
    }

    if true == appMgr.manageConfigMaps {
        appMgr.configMapLabel = cfgMapSelector.String()
        cfgMapOptions := func(options *metav1.ListOptions) {
            options.LabelSelector = appMgr.configMapLabel
        }
        log.Infof("[CORE] Watching ConfigMap resources.")
        syncInterval := resyncPeriod
        if appMgr.hubMode {
            syncInterval = hubModeInterval
        }
        appInf.cfgMapInformer = cache.NewSharedIndexInformer(
            cache.NewFilteredListWatchFromClient(
                appMgr.restClientv1,
                Configmaps,
                namespace,
                cfgMapOptions,
            ),
            &v1.ConfigMap{},
            syncInterval,
            cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
        )
    } else {
        log.Infof("[CORE] Not watching ConfigMap resources.")
    }

    if true == appMgr.manageIngress {
        log.Infof("[CORE] Handling Ingress resource events.")
        appInf.ingInformer.AddEventHandlerWithResyncPeriod(
            &cache.ResourceEventHandlerFuncs{
                AddFunc:    func(obj interface{}) { appMgr.enqueueIngress(obj, OprTypeCreate) },
                UpdateFunc: func(old, cur interface{}) { appMgr.enqueueIngress(cur, OprTypeUpdate) },
                DeleteFunc: func(obj interface{}) { appMgr.enqueueIngress(obj, OprTypeDelete) },
            },
            resyncPeriod,
        )
    } else {
        log.Infof("[CORE] Not handling Ingress resource events.")
    }

    if true == appMgr.manageConfigMaps {
        log.Infof("[CORE] Handling ConfigMap resource events.")
        syncInterval := resyncPeriod
        if appMgr.hubMode {
            syncInterval = hubModeInterval
        }
        appInf.cfgMapInformer.AddEventHandlerWithResyncPeriod(
            &cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) { appMgr.enqueueConfigMap(obj, OprTypeCreate) },
                UpdateFunc: func(old, cur interface{}) {
                    if appMgr.hubMode || !reflect.DeepEqual(old, cur) {
                        appMgr.enqueueConfigMap(cur, OprTypeUpdate)
                    }
                },
                DeleteFunc: func(obj interface{}) { appMgr.enqueueConfigMap(obj, OprTypeDelete) },
            },
            syncInterval,
        )

    } else {
        log.Infof("[CORE] Not handling ConfigMap resource events.")
    }

    appInf.svcInformer.AddEventHandlerWithResyncPeriod(
        &cache.ResourceEventHandlerFuncs{
            AddFunc:    func(obj interface{}) { appMgr.enqueueService(obj, OprTypeCreate) },
            UpdateFunc: func(old, cur interface{}) { appMgr.enqueueService(cur, OprTypeUpdate) },
            DeleteFunc: func(obj interface{}) { appMgr.enqueueService(obj, OprTypeDelete) },
        },
        resyncPeriod,
    )

    appInf.endptInformer.AddEventHandlerWithResyncPeriod(
        &cache.ResourceEventHandlerFuncs{
            AddFunc:    func(obj interface{}) { appMgr.enqueueEndpoints(obj, OprTypeCreate) },
            UpdateFunc: func(old, cur interface{}) { appMgr.enqueueEndpoints(cur, OprTypeUpdate) },
            DeleteFunc: func(obj interface{}) { appMgr.enqueueEndpoints(obj, OprTypeDelete) },
        },
        resyncPeriod,
    )

    return &appInf
}

func (appMgr *Manager) enqueueIngress(obj interface{}, operation string) {
    if ok, keys := appMgr.checkValidIngress(obj); ok {
        for _, key := range keys {
            key.Operation = operation
            appMgr.vsQueue.Add(*key)
            log.Infof("[CORE] Add %v to queue", key)
        }
    }
}

func (appMgr *Manager) enqueueConfigMap(obj interface{}, operation string) {
    if ok, keys := appMgr.checkValidConfigMap(obj, operation); ok {
        for _, key := range keys {
            key.Operation = operation
            appMgr.vsQueue.Add(*key)
            log.Infof("[CORE] Add %v to queue", key)
        }
    }
}

func (appMgr *Manager) enqueueService(obj interface{}, operation string) {
    if ok, keys := appMgr.checkValidService(obj); ok {
        for _, key := range keys {
            key.Operation = operation
            appMgr.vsQueue.Add(*key)
            log.Infof("[CORE] Add %v to queue", key)
        }
    }
}

func (appMgr *Manager) enqueueEndpoints(obj interface{}, operation string) {
    if ok, keys := appMgr.checkValidEndpoints(obj); ok {
        for _, key := range keys {
            key.Operation = operation
            appMgr.vsQueue.Add(*key)
            log.Infof("[CORE] Add %v to queue", key)
        }
    }
}

func (appMgr *Manager) syncNamespace(nsName string) error {
        startTime := time.Now()
        var err error
        defer func() {
                endTime := time.Now()
                log.Debugf("[CORE] Finished syncing namespace %+v (%v)",
                        nsName, endTime.Sub(startTime))
        }()
        _, exists, err := appMgr.nsInformer.GetIndexer().GetByKey(nsName)
        if nil != err {
                log.Warningf("[CORE] Error looking up namespace '%v': %v\n", nsName, err)
                return err
        }

        appMgr.informersMutex.Lock()
        defer appMgr.informersMutex.Unlock()
        appInf, found := appMgr.getNamespaceInformerLocked(nsName)
        if exists && found {
                appMgr.triggerSyncResources(nsName, appInf)
                return nil
        }
        // Skip deleting informers if watching specific namespaces
        if !appMgr.DynamicNS {
                return nil
        }

        if exists {
                // exists but not found in informers map, add
                cfgMapSelector, err := labels.Parse(DefaultConfigMapLabel)
                if err != nil {
                        return fmt.Errorf("Failed to parse Label Selector string: %v", err)
                }
                appInf, err = appMgr.addNamespaceLocked(nsName, cfgMapSelector, 0)
                if err != nil {
                        return fmt.Errorf("Failed to add informers for namespace %v: %v",
                                nsName, err)
                }
                appInf.start()
                appInf.waitForCacheSync()
        } else {
                // does not exist but found in informers map, delete
                // Clean up all resources that reference a removed namespace
                appInf.stopInformers()
                appMgr.removeNamespaceLocked(nsName)
                appMgr.eventNotifier.DeleteNotifierForNamespace(nsName)
                appMgr.resources.Lock()
                rsDeleted := 0
                appMgr.resources.ForEach(func(key ServiceKey, cfg *ResourceConfig) {
                        if key.Namespace == nsName {
                                if appMgr.resources.Delete(key, "") {
                                        rsDeleted += 1
                                }
                        }
                })
                appMgr.resources.Unlock()
                // Handle Agent Specific ConfigMaps
                if appMgr.AgentCIS.IsImplInAgent(ResourceTypeCfgMap) {
                        for _, cm := range appMgr.agentCfgMap {
                                if cm.Namespace == nsName {
                                        cm.Operation = OprTypeDelete
                                        rsDeleted += 1
                                }
                        }
                }
                if rsDeleted > 0 {
                        log.Warningf("[CORE] Error looking up namespace '%v': %v\n", nsName, err)
                        appMgr.deployResource()
                }
        }

        return nil
}

func (appMgr *Manager) triggerSyncResources(ns string, inf *appInformer) {
        enqueueSvcFromNamespace := func(namespace string, appInf *appInformer) {
                objs, err := appInf.svcInformer.GetIndexer().ByIndex("namespace", namespace)
                if err != nil {
                        log.Errorf("[CORE] Unable to fetch services from namespace: %v for periodic resync", namespace)
                        return
                }
                if objs != nil && len(objs) > 0 {
                        svc := objs[0].(*v1.Service)
                        svcKey := serviceQueueKey{
                                Namespace:    namespace,
                                ServiceName:  svc.Name,
                                ResourceKind: Services,
                                ResourceName: svc.Name,
                                Operation:    OprTypeUpdate,
                        }
                        log.Debugf("[CORE] Periodic enqueue of Service from Namespace: %v, svc: %s", namespace, svc.Name)
                        appMgr.vsQueue.Add(svcKey)
                }
        }

        if appMgr.watchingAllNamespacesLocked() {
                namespaces := appMgr.GetWatchedNamespacesLockless()

                if len(namespaces) == 1 && namespaces[0] == "" {
                        if appMgr.nsInformer == nil {
                                return
                        }
                        nsps := appMgr.nsInformer.GetIndexer().List()
                        namespaces = []string{}
                        for _, ns := range nsps {
                                namespaces = append(namespaces, ns.(*v1.Namespace).Name)
                        }
                }

                for _, ns := range namespaces {
                        if inf, ok := appMgr.getNamespaceInformerLocked(ns); ok {
                                enqueueSvcFromNamespace(ns, inf)
                        }
                }
        } else {
                enqueueSvcFromNamespace(ns, inf)
        }
}

func (appMgr *Manager) removeNamespaceLocked(namespace string) error {
        if _, found := appMgr.appInformers[namespace]; !found {
                return fmt.Errorf("No informers exist for namespace %v\n", namespace)
        }
        delete(appMgr.appInformers, namespace)
        return nil
}

func (appMgr *Manager) GetWatchedNamespacesLockless() []string {
        var namespaces []string
        for k, _ := range appMgr.appInformers {
                namespaces = append(namespaces, k)
        }
        return namespaces
}


func (appMgr *Manager) syncVirtualServer(sKey serviceQueueKey) error {

    startTime := time.Now()
    defer func() {
        endTime := time.Now()
                // processedItems with +1 because that is the actual number of items processed
                // and it gets incremented just after this function returns
         log.Debugf("[CORE] Finished syncing virtual servers %+v in namespace %+v (%v), %v/%v", sKey.ServiceName, sKey.Namespace, endTime.Sub(startTime), appMgr.processedItems+1, appMgr.queueLen)
    }()
        // Get the informers for the namespace. This will tell us if we care about
        // this item.
    appInf, haveNamespace := appMgr.getNamespaceInformer(sKey.Namespace)
    if !haveNamespace {
                // This shouldn't happen as the namespace is checked for every item before
                // it is added to the queue, but issue a warning if it does.
        log.Warningf("Received an update for an item from an un-watched namespace %v", sKey.Namespace)
        return nil
    }

     // Lookup the service
    svcKey := sKey.Namespace + "/" + sKey.ServiceName
    obj, svcFound, err := appInf.svcInformer.GetIndexer().GetByKey(svcKey)
    if nil != err {
                // Returning non-nil err will re-queue this item with rate-limiting.
        log.Warningf("[CORE] Error looking up service '%v': %v\n", svcKey, err)
        return err
    }

    // Processing just one service from a namespace processes all the resources in that namespace
    switch sKey.ResourceKind {
    case Services:
        rkey := Services + "_" + sKey.Namespace
        if !appMgr.steadyState && sKey.Operation == OprTypeCreate {
            if _, ok := appMgr.processedResources[rkey]; ok {
                if !appMgr.steadyState && appMgr.processedItems >= appMgr.queueLen-1 {
                    appMgr.deployResource()
                    appMgr.steadyState = true
                }
                return nil
            }
            appMgr.processedResourcesMutex.Lock()
            appMgr.processedResources[rkey] = true
            appMgr.processedResourcesMutex.Unlock()
        }
    case Endpoints:
        if appMgr.IsNodePort() {
            return nil
        }
    case Configmaps:
        resKey := prepareResourceKey(sKey.ResourceKind, sKey.Namespace, sKey.ResourceName)
        switch sKey.Operation {
        case OprTypeCreate:
            if _, ok := appMgr.processedResources[resKey]; ok {
                if !appMgr.steadyState && appMgr.processedItems >= appMgr.queueLen-1 {
                    appMgr.deployResource()
                    appMgr.steadyState = true
                }
                return nil
            }
        case OprTypeDelete:
            appMgr.processedResourcesMutex.Lock()
            delete(appMgr.processedResources, resKey)
            appMgr.processedResourcesMutex.Unlock()
        }
    default:
                // Resources other than Services will be tracked if they are processed earlier
        resKey := prepareResourceKey(sKey.ResourceKind, sKey.Namespace, sKey.ResourceName)
        switch sKey.Operation {
                // If a resource is processed earlier and still sKey gives us CREATE event,
                // then it was handled earlier when associated service processed
                // otherwise just mark it as processed and continue
        case OprTypeCreate:
            if _, ok := appMgr.processedResources[resKey]; ok {
                if !appMgr.steadyState && appMgr.processedItems >= appMgr.queueLen-1 {
                    appMgr.deployResource()
                    appMgr.steadyState = true
                }
                return nil
            }
            appMgr.processedResourcesMutex.Lock()
            appMgr.processedResources[resKey] = true
            appMgr.processedResourcesMutex.Unlock()
        case OprTypeDelete:
            appMgr.processedResourcesMutex.Lock()
            delete(appMgr.processedResources, resKey)
            appMgr.processedResourcesMutex.Unlock()
        }

    }

    // Use a map to allow ports in the service to be looked up quickly while
    // looping through the ConfigMaps. The value is not currently used.
    svcPortMap := make(map[int32]bool)
    var svc *v1.Service
    if svcFound {
        svc = obj.(*v1.Service)
        for _, portSpec := range svc.Spec.Ports {
            svcPortMap[portSpec.Port] = false
        }
    }

    // rsMap stores all resources currently in Resources matching sKey, indexed by port.
    // At the end of processing, rsMap should only contain configs we want to delete.
    // If we have a valid config, then we remove it from rsMap.
    rsMap := appMgr.getResourcesForKey(sKey)
    dgMap := make(InternalDataGroupMap)

    var stats vsSyncStats
    appMgr.rsrcSSLCtxt = make(map[string]*v1.Secret)
    if nil != appInf.ingInformer {
        err = appMgr.syncIngresses(&stats, sKey, rsMap, svcPortMap, svc, appInf, dgMap)
        if nil != err {
            return err
        }
    }
    if nil != appInf.cfgMapInformer {
        err = appMgr.syncConfigMaps(&stats, sKey, rsMap, svcPortMap, svc, appInf)
        if nil != err {
            return err
        }
    }
        // Update internal data groups if changed
    appMgr.syncDataGroups(&stats, dgMap, sKey.Namespace)
        // Delete IRules if necessary
    appMgr.syncIRules()

    if len(rsMap) > 0 {
                // We get here when there are ports defined in the service that don't
                // have a corresponding config map.
        stats.vsDeleted += appMgr.deleteUnusedConfigs(sKey, rsMap)
        stats.vsUpdated += appMgr.deleteUnusedResources(sKey, svcFound)

    } else if !svcFound {
        stats.vsUpdated += appMgr.deleteUnusedResources(sKey, svcFound)
    }

    log.Debugf("[CORE] Updated %v of %v virtual server configs, deleted %v", stats.vsUpdated, stats.vsFound, stats.vsDeleted)

        // delete any custom profiles that are no longer referenced
    appMgr.deleteUnusedProfiles(appInf, sKey.Namespace, &stats)

    switch {
    case stats.isStatsAvailable(),
        !appMgr.steadyState && appMgr.processedItems >= appMgr.queueLen-1:
                {
                        if appMgr.processedItems >= appMgr.queueLen-1 || appMgr.steadyState {
                                appMgr.deployResource()
                                appMgr.steadyState = true
                        }
                }
    }

    return nil
}

func prepareResourceKey(kind, namespace, name string) string {
    return kind + "_" + namespace + "/" + name
}

func (appMgr *Manager) getResourcesForKey(sKey serviceQueueKey) ResourceMap {
        // Return a copy of what is stored in resources
    appMgr.resources.Lock()
    defer appMgr.resources.Unlock()
    rsMap := make(ResourceMap)
    appMgr.resources.ForEach(func(key ServiceKey, cfg *ResourceConfig) {
        if key.Namespace == sKey.Namespace && key.ServiceName == sKey.ServiceName {
            rsMap[key.ServicePort] = append(rsMap[key.ServicePort], cfg)
        }
    })
    return rsMap
}

func (appMgr *Manager) deleteUnusedResources(sKey serviceQueueKey, svcFound bool,) int {
    appMgr.resources.Lock()
    defer appMgr.resources.Unlock()
    rsUpdated := 0
    namespace := sKey.Namespace
    svcName := sKey.ServiceName
    for _, cfg := range appMgr.resources.GetAllResources() {
        if cfg.MetaData.ResourceType == "configmap" ||
            cfg.MetaData.ResourceType == "iapp" {
            continue
        }
        for _, pool := range cfg.Pools {
                        // Make sure we aren't processing empty pool
            if pool.Name != "" {
                key := ServiceKey{
                    ServiceName: pool.ServiceName,
                    ServicePort: pool.ServicePort,
                    Namespace:   namespace,
                }
                poolNS := strings.Split(pool.Name, "_")[1]
                _, ok := appMgr.resources.Get(key, cfg.GetName())
                if pool.ServiceName == svcName && poolNS == namespace && (!ok || !svcFound) {
                    if updated, svcKey := cfg.RemovePool(namespace, pool.Name, appMgr.mergedRulesMap); updated {
                        appMgr.resources.DeleteKeyRefLocked(*svcKey, cfg.GetName())
                        rsUpdated += 1
                    }
                }
            }
        }
    }
    return rsUpdated
}

func (appMgr *Manager) deleteUnusedConfigs(
        sKey serviceQueueKey,
        rsMap ResourceMap,
) int {
        rsDeleted := 0
        appMgr.resources.Lock()
        defer appMgr.resources.Unlock()
        // First delete any configs that we have left over from processing
        // (Configs that are still valid aren't left over)
        for port, cfgList := range rsMap {
                tmpKey := ServiceKey{
                        Namespace:   sKey.Namespace,
                        ServiceName: sKey.ServiceName,
                        ServicePort: port,
                }
                for _, cfg := range cfgList {
                        rsName := cfg.GetName()
                        if appMgr.resources.Delete(tmpKey, rsName) {
                                rsDeleted += 1
                        }
                }
        }
        return rsDeleted
}
