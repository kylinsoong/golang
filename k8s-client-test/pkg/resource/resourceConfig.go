package resource

import (
    "sync"
)

type Resources struct {
    sync.Mutex
    rm      resourceKeyMap
    RsMap   ResourceConfigMap
    objDeps ObjectDependencyMap
}

type ObjectDependency struct {
    Kind      string
    Namespace string
    Name      string
}

type ObjectDependencies map[ObjectDependency]int

type resourceList map[string]bool

type resourceKeyMap map[ServiceKey]resourceList

type ResourceConfigMap map[string]*ResourceConfig

type ObjectDependencyMap map[ObjectDependency]ObjectDependencies

type ResourceEnumFunc func(key ServiceKey, cfg *ResourceConfig)

func NewResources() *Resources {
    var rs Resources
    rs.Init()
    return &rs
}

func (rs *Resources) Init() {
    rs.rm = make(resourceKeyMap)
    rs.RsMap = make(ResourceConfigMap)
    rs.objDeps = make(ObjectDependencyMap)
}

func (rs *Resources) ForEach(f ResourceEnumFunc) {
    for svcKey, rsList := range rs.rm {
        for rsName, _ := range rsList {
            cfg, _ := rs.RsMap[rsName]
            f(svcKey, cfg)
        }
    }
}

func (rs *Resources) deleteImpl(
        rsList resourceList,
        rsName string,
        svcKey ServiceKey,
) {
        //bigIPPrometheus.MonitoredServices.DeleteLabelValues(svcKey.Namespace, svcKey.ServiceName, "parse-error")
        //bigIPPrometheus.MonitoredServices.DeleteLabelValues(svcKey.Namespace, rsName, "port-not-found")
        //bigIPPrometheus.MonitoredServices.DeleteLabelValues(svcKey.Namespace, rsName, "service-not-found")
        //bigIPPrometheus.MonitoredServices.DeleteLabelValues(svcKey.Namespace, rsName, "success")

        // Remove mapping for a backend -> virtual/iapp
        delete(rsList, rsName)
        if len(rsList) == 0 {
                // Remove backend since no virtuals/iapps remain
                delete(rs.rm, svcKey)
        }

        // Look at all service keys to see if another references rsName
        useCt := 0
        for _, otherList := range rs.rm {
                for otherName := range otherList {
                        if otherName == rsName {
                                // Found one, can't delete this resource yet.
                                useCt += 1
                                break
                        }
                }
        }
        if useCt == 0 {
                delete(rs.RsMap, rsName)
        }
}

func (rs *Resources) Delete(svcKey ServiceKey, name string) bool {
        rsList, ok := rs.rm[svcKey]
        if !ok {
                // svcKey not found
                return false
        }
        if name == "" {
                // Delete all resources for svcKey
                for rsName, _ := range rsList {
                        rs.deleteImpl(rsList, rsName, svcKey)
                }
                return true
        }
        if _, ok = rsList[name]; ok {
                // Delete specific named resource for svcKey
                rs.deleteImpl(rsList, name, svcKey)
                return true
        }
        return false
}

