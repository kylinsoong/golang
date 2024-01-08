package agent

import (
    . "github.com/kylinsoong/bigip-ctlr/pkg/agent/as3"
    "github.com/kylinsoong/bigip-ctlr/pkg/resource"
)

type agentAS3 struct {
    *AS3Manager
}

func (ag *agentAS3) Init(params interface{}) error {

    return nil
}

func (ag *agentAS3) GetBigipRegKey() string {

    return "key"
}

func (ag *agentAS3) Deploy(req interface{}) error {

    return nil
}

func (ag *agentAS3) Clean(partition string) error {

    return nil
}

func (ag *agentAS3) DeInit() error {

    return nil
}

func (ag *agentAS3) IsImplInAgent(rsrc string) bool {
    if resource.ResourceTypeCfgMap == rsrc {
        return true
    }
    return false
}

