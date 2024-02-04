package agent

import (
    log  "github.com/kylinsoong/golang/nim-rest-client/pkg/vlogger"
)

type agentAS3 struct {
}

func (ag *agentAS3) Init(params interface{}) error {
    log.Info("[AS3] Initializing AS3 Agent")
    return nil
}

func (ag *agentAS3) DoBatchDeploy(groups, stage string) error {
    return nil
}

func (ag *agentAS3) DoDeploy(instance, stage string) error {
    return nil
}

func (ag *agentAS3) Clean(partition string) error {
    return nil
}

func (ag *agentAS3) DeInit() error {
    return nil
}

func (ag *agentAS3) IsImplInAgent(rsrc string) bool {
    return false
}


