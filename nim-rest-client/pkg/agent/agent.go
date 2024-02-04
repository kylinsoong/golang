package agent

import (
    "errors"
)

const (
    AS3Agent = "as3"
    NIMAgent = "nim"
)

type RESTAgentInterface interface {
    Initializer
    Deployer
    DeInitializer
    Remover
    IsImplInAgent(string) bool
}

type Initializer interface {
    Init(interface{}) error
}

type Deployer interface {
    DoBatchDeploy(groups, stage string) error
    DoDeploy(instance, stage string) error
}

type DeInitializer interface {
    DeInit() error
}

type Remover interface {
    Clean(partition string) error
}

func CreateAgent(agentType string) (RESTAgentInterface, error) {
        switch agentType {
        case AS3Agent:
                return new(agentAS3), nil
        case NIMAgent:
                return new(agentNIM), nil
        default:
                return nil, errors.New("Invalid Agent Type")
        }
}
