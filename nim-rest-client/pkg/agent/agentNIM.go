package agent

import (
    . "github.com/kylinsoong/golang/nim-rest-client/pkg/agent/nim"
    log  "github.com/kylinsoong/golang/nim-rest-client/pkg/vlogger"
)

type agentNIM struct {
    *NIMManager
}

func (ag *agentNIM) Init(params interface{}) error {
    nimParams := params.(*Params)
    log.Infof("[NIM] Initializing NIM REST Agent, %v", nimParams)
    ag.NIMManager = NewNIMManager(nimParams)
    return nil
}

func (ag *agentNIM) DoBatchDeploy(groups, stage string) error {
    err := ag.RestUpdateInstanceGroup(groups, stage)
    if err != nil {
        return err
    }
    return nil
}

func (ag *agentNIM) DoDeploy(instance, stage string) error {
    err := ag.RestUpdateInstance(instance, stage)
    if err != nil {
        return err
    }
    return nil
}


func (ag *agentNIM) Clean(partition string) error {
    return nil
}

func (ag *agentNIM) DeInit() error {
    return nil
}

func (ag *agentNIM) IsImplInAgent(rsrc string) bool {
    return true
}


