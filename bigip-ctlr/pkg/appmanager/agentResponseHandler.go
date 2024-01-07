package appmanager

import (
    log "github.com/kylinsoong/bigip-ctlr/pkg/vlogger"
)

func (appMgr *Manager) agentResponseWorker() {
    log.Debugf("[CORE] Agent Response Worker started and blocked on channel %v", appMgr.agRspChan)
}
