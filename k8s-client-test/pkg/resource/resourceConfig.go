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
