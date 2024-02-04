package nim

import (
    //"github.com/kylinsoong/golang/nim-rest-client/pkg/resources"
    log  "github.com/kylinsoong/golang/nim-rest-client/pkg/vlogger"
)

type NIMManager struct {
    sslInsecure         bool
    enableTLS           string
    ciphers             string
    userAgent           string
    PostManager         *PostManager
}

type Params struct {
    EnableTLS           string
    TrustedCerts        string
    Ciphers             string
    UserAgent           string
    NIMUsername         string
    NIMPassword         string
    NIMURL              string
    SSLInsecure         bool
}

func NewNIMManager(params *Params) *NIMManager{
    nimManager := NIMManager{
        userAgent:                 params.UserAgent,
        sslInsecure:               params.SSLInsecure,
	enableTLS:                 params.EnableTLS,
        ciphers:                   params.Ciphers,
        PostManager: NewPostManager(PostParams{
	    NIMUsername:   params.NIMUsername,
	    NIMPassword:   params.NIMPassword,
	    NIMURL:        params.NIMURL,
	    TrustedCerts:  params.TrustedCerts,
	    SSLInsecure:   params.SSLInsecure,}),
    }

    return &nimManager
}

func (nm *NIMManager) RestUpdateInstanceGroup(group, stage string) error {

    //var instanceGroupConfigResponse resources.InstanceGroupConfigResponse

    if stage != "" && len(stage) > 1 {
        log.Debugf("[NIM] API Process Against InstanceGroup %s, with stage config %s enabled", group, stage)
        group_uid, err := nm.PostManager.GetInstanceGroupUid(group)
        if err != nil {
            return err
        }
        log.Infof("[NIM] %s referenced group uid is %s", group, group_uid)
    } else {
        log.Debugf("[NIM] API Process Against InstanceGroup %s", group)
        group_uid, err := nm.PostManager.GetInstanceGroupUid(group)
        if err != nil {
            return err
        }
        log.Infof("[NIM] %s referenced group uid is %s", group, group_uid)
       
        instanceGroupConfigResponse, err := nm.PostManager.GetInstanceGroupConfig(group_uid)
        if err != nil {
            return err
        }

        log.Infof("%v", instanceGroupConfigResponse)
    } 
    
    return nil
}

func (nm *NIMManager) RestUpdateInstance(instance, stage string) error {

    if stage != "" && len(stage) > 1 {
        log.Debugf("[NIM] API Process Against Instance %s, with stage config %s enabled", instance, stage)
    } else {
        log.Debugf("[NIM] API Process Against Instance %s", instance)
    }

    return nil
}
