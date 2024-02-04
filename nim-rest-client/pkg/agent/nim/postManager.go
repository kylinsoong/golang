package nim

import (
    "crypto/x509"
    "net/http"
    "crypto/tls"
    "time"
    "encoding/json"
    "fmt"
    "strings"
    "io/ioutil"
    . "github.com/kylinsoong/golang/nim-rest-client/pkg/resources"
    log  "github.com/kylinsoong/golang/nim-rest-client/pkg/vlogger"
)

const (
	timeoutNill   = 0 * time.Second
	timeoutSmall  = 3 * time.Second
	timeoutMedium = 30 * time.Second
	timeoutLarge  = 60 * time.Second
)

type PostManager struct {
    httpClient *http.Client
    PostParams
}

type PostParams struct {
    NIMUsername     string
    NIMPassword     string
    NIMURL          string
    TrustedCerts    string
    SSLInsecure     bool
}

func NewPostManager(params PostParams) *PostManager {
    pm := &PostManager{
        PostParams: params,
    }
    pm.setupNIMRESTClient()

    return pm
}

func (postMgr *PostManager) setupNIMRESTClient() {

    rootCAs, _ := x509.SystemCertPool()
    if rootCAs == nil {
	rootCAs = x509.NewCertPool()
    }

    certs := []byte(postMgr.TrustedCerts)
    if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
	log.Debug("[AS3] No certs appended, using only system certs")
    }

    tr := &http.Transport{
	TLSClientConfig: &tls.Config{
			InsecureSkipVerify: postMgr.SSLInsecure,
			RootCAs:            rootCAs,
	},
    }

    postMgr.httpClient = &http.Client{
	Transport: tr,
	Timeout:   timeoutLarge,
    }

} 

func (postMgr *PostManager) httpReq(request *http.Request) (*http.Response, interface{}) {
	httpResp, err := postMgr.httpClient.Do(request)
	if err != nil {
		log.Errorf("[NIM] REST call error: %v ", err)
		return nil, nil
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorf("[NIM] REST call response error: %v ", err)
		return nil, nil
	}
	return httpResp, body
}

func (postMgr *PostManager) getInstanceGroupSummaryURL() string {
    apiURL := postMgr.NIMURL + "/api/platform/v1/instance-groups/summary"
    return apiURL
}

func (postMgr *PostManager) GetInstanceGroupUid(group string) (string, error) {

    url := postMgr.getInstanceGroupSummaryURL()
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
	log.Errorf("[NIM] Creating new HTTP request error: %v ", err)
	return "", err
    }
 
    log.Infof("[NIM] posting GET NIM Instance Group Summary on %v", url)
    req.SetBasicAuth(postMgr.NIMUsername, postMgr.NIMPassword)

    httpResp, response := postMgr.httpReq(req)
    if httpResp == nil || response == nil {
	return "", fmt.Errorf("Internal Error")
    }

    bodyBytes, ok := response.([]byte)
    if !ok {
        return "", fmt.Errorf("Failed to convert response to []byte")
    }

    var loadInstanceGroups LoadInstanceGroups
    err = json.Unmarshal(bodyBytes, &loadInstanceGroups)
    log.Debugf("[NIM] HTTP Response: %v", loadInstanceGroups)
    if err != nil {
	return "", fmt.Errorf("Error unmarshaling JSON: %v", err)
    }

    for _, item := range loadInstanceGroups.Items{
        if item.Name == group {
            return item.UID, nil
        }
    }
   
    return "", fmt.Errorf("can not find group uid via %s", group)
} 

func (postMgr *PostManager) getInstanceGroupConfigURL(uid string) string {
    apiURL := strings.Replace("/api/platform/v1/instance-groups/{uid}/config", "{uid}", uid, 1)
    return postMgr.NIMURL + apiURL
}

func (postMgr *PostManager) GetInstanceGroupConfig(uid string) (*InstanceGroupConfigResponse, error) {

    url := postMgr.getInstanceGroupConfigURL(uid)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Errorf("[NIM] Creating new HTTP request error: %v ", err)
        return nil, err
    }

    log.Infof("[NIM] posting GET NIM Instance Group Config on %v", url)
    req.SetBasicAuth(postMgr.NIMUsername, postMgr.NIMPassword)

    httpResp, response := postMgr.httpReq(req)
    if httpResp == nil || response == nil {
        return nil, fmt.Errorf("Internal Error")
    }

    bodyBytes, ok := response.([]byte)
    if !ok {
        return nil, fmt.Errorf("Failed to convert response to []byte")
    }

    var instanceGroupConfigResponse InstanceGroupConfigResponse
    err = json.Unmarshal(bodyBytes, &instanceGroupConfigResponse)
    log.Debugf("[NIM] HTTP Response: %v", instanceGroupConfigResponse)
    if err != nil {
        return nil, fmt.Errorf("Error unmarshaling JSON: %v", err)
    }

    return &instanceGroupConfigResponse, nil

}
