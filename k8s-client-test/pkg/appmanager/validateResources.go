package appmanager

import (
    v1 "k8s.io/api/core/v1"
    //. "github.com/kylinsoong/k8s-client-test/pkg/resource"
)

func (appMgr *Manager) checkValidIngress(obj interface{},) (bool, []*serviceQueueKey) {
    return false, nil
}

func (appMgr *Manager) checkValidConfigMap(obj interface{}, oprType string,) (bool, []*serviceQueueKey) {
    var keyList []*serviceQueueKey
    cm := obj.(*v1.ConfigMap)
    namespace := cm.ObjectMeta.Namespace
    _, ok := appMgr.getNamespaceInformer(namespace)
    if !ok {
                // Not watching this namespace
        return false, nil
    }
    key := &serviceQueueKey{
        Namespace:    namespace,
        Operation:    oprType,
        ResourceKind: Configmaps,
        ResourceName: cm.Name,
    }
    keyList = append(keyList, key)
    return true, keyList
}

func (appMgr *Manager) checkValidService(obj interface{},) (bool, []*serviceQueueKey) {

    svc := obj.(*v1.Service)
    namespace := svc.ObjectMeta.Namespace
    _, ok := appMgr.getNamespaceInformer(namespace)
    if !ok {
                // Not watching this namespace
        return false, nil
    }
    key := &serviceQueueKey{
        ServiceName:  svc.ObjectMeta.Name,
        Namespace:    namespace,
        ResourceKind: Services,
        ResourceName: svc.Name,
    }
    var keyList []*serviceQueueKey
    keyList = append(keyList, key)
    return true, keyList
}

func (appMgr *Manager) checkValidEndpoints(obj interface{},) (bool, []*serviceQueueKey) {
    eps := obj.(*v1.Endpoints)
    namespace := eps.ObjectMeta.Namespace
        // Check if the service to see if we care about it.
    _, ok := appMgr.getNamespaceInformer(namespace)
    if !ok {
            // Not watching this namespace
            return false, nil
    }
    key := &serviceQueueKey{
        ServiceName:  eps.ObjectMeta.Name,
        Namespace:    namespace,
        ResourceKind: Endpoints,
        ResourceName: eps.Name,
    }
    var keyList []*serviceQueueKey
    keyList = append(keyList, key)
    return true, keyList
}

