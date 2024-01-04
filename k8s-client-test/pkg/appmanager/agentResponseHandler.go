package appmanager

import (
    log "github.com/kylinsoong/k8s-client-test/pkg/vlogger"
)

func (appMgr *Manager) agentResponseWorker() {
    log.Debugf("[CORE] Agent Response Worker started and blocked on channel %v", appMgr.agRspChan)
}
