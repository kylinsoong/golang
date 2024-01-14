/*
This test is simulate the as3 Agent behavior.

1. as3 Agent initialization:
 1. TestIsBigIPAppServicesAvailable()
 2. TestPartitionClean()
 3. TestGetBigipRegKey()
*/
package agent

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	cisAgent "github.com/kylinsoong/bigip-ctlr/pkg/agent"
	"github.com/kylinsoong/bigip-ctlr/pkg/appmanager"
	. "github.com/kylinsoong/bigip-ctlr/pkg/resource"
	"github.com/kylinsoong/bigip-ctlr/test"
)

var (
	appMgr     appmanager.Manager
	httpClient *http.Client
	globalMap  map[string]bool
)

type (
	as3Template        string
	as3Declaration     string
	as3ADC             as3JSONWithArbKeys
	as3Control         as3JSONWithArbKeys
	as3Tenant          as3JSONWithArbKeys
	as3Application     as3JSONWithArbKeys
	as3JSONWithArbKeys map[string]interface{}

	poolName   string
	appName    string
	tenantName string
	tenant     map[appName][]poolName
	as3Object  map[tenantName]tenant
)

type config struct {
	data      string
	as3APIURL string
}

func generateIP() string {
	ip := fmt.Sprintf("10.244.%d.%d", rand.Intn(256), rand.Intn(255))
	if globalMap[ip] {
		return generateIP()
	} else {
		globalMap[ip] = true
		return ip
	}
}

func post_as3_declaration(cfg config) {
	httpReqBody := bytes.NewBuffer([]byte(cfg.data))

	req, err := http.NewRequest("POST", cfg.as3APIURL, httpReqBody)
	if err != nil {
		fmt.Errorf("[AS3] Creating new HTTP request error: %v ", err)
	}
	fmt.Printf("[AS3] posting request to %v\n", cfg.as3APIURL)

	req.SetBasicAuth("admin", "admin")
	httpClient = createHTTPClient()
	httpResp, responseMap := httpReq(req)
	fmt.Printf("httpResp: %v\n", httpResp)
	fmt.Printf("response: %v\n", responseMap)

	results := (responseMap["results"]).([]interface{})
	for _, value := range results {
		v := value.(map[string]interface{})
		fmt.Printf("[AS3] Response from BIG-IP: code: %v, tenant: %v, message: %v, runTime: %v", v["code"], v["tenant"], v["message"], v["runTime"])
	}
}

func process_deletion(partition string) {
	emptyAS3Declaration := getEmptyAs3Declaration(partition)
	data := string(emptyAS3Declaration)
	url := getAS3APIURL([]string{partition})

	fmt.Printf("declaration: %s\n", data)
	fmt.Printf("url: %s\n", url)

	cfg := config{
		data:      data,
		as3APIURL: url,
	}
	httpReqBody := bytes.NewBuffer([]byte(cfg.data))

	req, err := http.NewRequest("POST", cfg.as3APIURL, httpReqBody)
	if err != nil {
		fmt.Errorf("[AS3] Creating new HTTP request error: %v ", err)
	}
	fmt.Printf("[AS3] posting request to %v\n", cfg.as3APIURL)

	req.SetBasicAuth("admin", "admin")
	httpClient = createHTTPClient()
	httpResp, responseMap := httpReq(req)
	fmt.Printf("httpResp: %v\n", httpResp)
	fmt.Printf("response: %v\n", responseMap)

	results := (responseMap["results"]).([]interface{})
	for _, value := range results {
		v := value.(map[string]interface{})
		fmt.Printf("[AS3] Response from BIG-IP: code: %v, tenant: %v, message: %v, runTime: %v\n", v["code"], v["tenant"], v["message"], v["runTime"])
	}
}

func exreact_tenant_names(name string) []string {
	tenantMap, _ := processCfgMap(name)
	partitions := make([]string, 0, len(tenantMap))
	for partition := range tenantMap {
		partitions = append(partitions, partition)
	}
	return partitions
}

func createHTTPClient() *http.Client {
	rootCAs, _ := x509.SystemCertPool()
	certs := []byte("")
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {

	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            rootCAs,
		},
	}

	timeoutLarge := 60 * time.Second
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   timeoutLarge,
	}
	return httpClient
}

func httpReq(request *http.Request) (*http.Response, map[string]interface{}) {
	httpResp, err := httpClient.Do(request)
	if err != nil {
		fmt.Printf("[AS3] REST call error: %v ", err)
		os.Exit(1)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Printf("[AS3] REST call response error: %v ", err)
		os.Exit(1)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("[AS3] Response body unmarshal failed: %v\n", err)
		os.Exit(1)
	}
	return httpResp, response
}

func TestCreateAgent(t *testing.T) {
	var err error
	agent := "as3"
	appMgr.AgentCIS, err = cisAgent.CreateAgent(agent)
	if err != nil {
		t.Logf("[INIT] unable to create agent %v error: err: %+v", agent, err)
		os.Exit(1)
	}
	t.Logf("type of : %s", reflect.TypeOf(appMgr.AgentCIS))
}

/*
The main AgentCIS.Init() logic is to verify the whether BIG-IP AS3 is available via execute the

	https://192.168.45.52/mgmt/shared/appsvcs/info

and check the respone body, and extract the as3 version, if the version is expected version, then the BIG-IP AS3 is available.

# Use the curl to simulate the process

% curl -s -k -u "admin:admin" -X GET https://192.168.45.52/mgmt/shared/appsvcs/info

	{
	  "version": "3.36.1",
	  "release": "1",
	  "schemaCurrent": "3.36.0",
	  "schemaMinimum": "3.0.0"
	}
*/
func TestIsBigIPAppServicesAvailable(t *testing.T) {

	url := "https://192.168.45.52/mgmt/shared/appsvcs/info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Logf("[AS3] Creating new HTTP request error: %v ", err)
		os.Exit(1)
	}
	t.Logf("[AS3] posting GET BIGIP AS3 Version request on %v", url)

	req.SetBasicAuth("admin", "admin")

	httpClient = createHTTPClient()

	httpResp, responseMap := httpReq(req)

	t.Logf("httpResp: %v", httpResp)
	t.Logf("response: %v", responseMap)

	as3VersionStr := responseMap["version"].(string)
	as3versionreleaseStr := responseMap["release"].(string)
	as3SchemaVersion := responseMap["schemaCurrent"].(string)
	t.Logf("as3VersionStr: %s", as3VersionStr)
	t.Logf("as3versionreleaseStr: %s", as3versionreleaseStr)
	t.Logf("as3SchemaVersion: %s", as3SchemaVersion)

}

var baseAS3Config = `{
	"$schema": "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/master/schema/%s/as3-schema-%s.json",
	"class": "AS3",
	"declaration": {
	  "class": "ADC",
	  "schemaVersion": "%s",
	  "id": "urn:uuid:85626792-9ee7-46bb-8fc8-4ba708cfdc1d",
	  "label": "CIS Declaration",
	  "remark": "Auto-generated by CIS",
	  "controls": {
		 "class": "Controls",
		 "userAgent": "CIS Configured AS3"
	  }
	}
  }
`

/*
This function will create a empty as3 declaration file:

	{
	  "$schema": "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/master/schema/3.36.1/as3-schema-3.36.1-1.json",
	  "class": "AS3",
	  "declaration": {
	    "class": "ADC",
	    "controls": {
	      "class": "Controls",
	      "userAgent": "CIS/v K8S/v 1.23.10"
	    },
	    "id": "urn:uuid:85626792-9ee7-46bb-8fc8-4ba708cfdc1d",
	    "k8s": {
	      "Shared": {
	        "class": "Application",
	        "template": "shared"
	      },
	      "class": "Tenant",
	      "defaultRouteDomain": 0
	    },
	    "label": "CIS Declaration",
	    "remark": "Auto-generated by CIS",
	    "schemaVersion": "3.36.0"
	  }
	}
*/
func createEmptyAS3Declaration() as3Declaration {
	var as3Config map[string]interface{}
	baseAS3ConfigEmpty := fmt.Sprintf(baseAS3Config, "3.36.1", "3.36.1-1", "3.36.0")
	_ = json.Unmarshal([]byte(baseAS3ConfigEmpty), &as3Config)
	decl := as3Config["declaration"].(map[string]interface{})
	controlObj := make(as3Control)
	controlObj["class"] = "Controls"
	controlObj["userAgent"] = "CIS/v K8S/v 1.23.10"
	decl["controls"] = controlObj
	partition := "k8s"
	tenantObj := make(as3Tenant)
	app := as3Application{}
	app["class"] = "Application"
	app["template"] = "shared"
	tenantObj["class"] = "Tenant"
	tenantObj["Shared"] = app
	tenantObj["defaultRouteDomain"] = 0
	decl[partition] = tenantObj
	data, _ := json.Marshal(as3Config)
	emptyAS3Declaration := as3Declaration(data)
	return emptyAS3Declaration
}

/*
The steps of partition clean:
 1. create a empty as3 declaration
 2. create HTTP POST request uss empty as3 declaration as body, and "/mgmt/shared/appsvcs/declare/k8s" as url
 3. create HTTP Client, execute HTTP POST Request
 4. verify the response, print response message

if use curl to simulate above steps it looks below:

% curl -s -k -u "admin:admin" -X POST -H "Content-Type: application/json" -d @$(pwd)/emptyAS3Declaration.json  https://192.168.45.52/mgmt/shared/appsvcs/declare/k8s

	{
	  "results": [
	    {
	      "code": 200,
	      "message": "no change",
	      "host": "localhost",
	      "tenant": "k8s",
	      "runTime": 893
	    }
	  ],
	  "declaration": {
	    "k8s": {
	      "Shared": {
	        "class": "Application",
	        "template": "shared"
	      },
	      "class": "Tenant",
	      "defaultRouteDomain": 0
	    },
	    "class": "ADC",
	    "controls": {
	      "class": "Controls",
	      "userAgent": "CIS/v K8S/v 1.23.10",
	      "archiveTimestamp": "2024-01-13T02:06:44.596Z"
	    },
	    "id": "urn:uuid:85626792-9ee7-46bb-8fc8-4ba708cfdc1d",
	    "label": "CIS Declaration",
	    "remark": "Auto-generated by CIS",
	    "schemaVersion": "3.36.0",
	    "updateMode": "selective"
	  }
	}
*/
func TestPartitionClean(t *testing.T) {

	emptyAS3Declaration := createEmptyAS3Declaration()
	t.Logf("emptyAS3Declaration: %v", emptyAS3Declaration)

	data := string(emptyAS3Declaration)
	url := "https://192.168.45.52/mgmt/shared/appsvcs/declare/k8s"
	cfg := config{
		data:      data,
		as3APIURL: url,
	}
	httpReqBody := bytes.NewBuffer([]byte(cfg.data))

	req, err := http.NewRequest("POST", cfg.as3APIURL, httpReqBody)
	if err != nil {
		t.Errorf("[AS3] Creating new HTTP request error: %v ", err)
	}
	t.Logf("[AS3] posting request to %v", cfg.as3APIURL)

	req.SetBasicAuth("admin", "admin")
	httpClient = createHTTPClient()
	httpResp, responseMap := httpReq(req)
	t.Logf("httpResp: %v", httpResp)
	t.Logf("response: %v", responseMap)

	results := (responseMap["results"]).([]interface{})
	for _, value := range results {
		v := value.(map[string]interface{})
		t.Logf("[AS3] Response from BIG-IP: code: %v, tenant: %v, message: %v, runTime: %v", v["code"], v["tenant"], v["message"], v["runTime"])
	}
}

/*
This steps is execute the HTTP GET to extract license key, it's sammiliar like the below curl

% curl -s -k -u "admin:admin" -X GET https://192.168.45.52/mgmt/tm/shared/licensing/registration | jq .registrationKey
"KVPKO-EBYPF-UFQQG-WYBNP-TXRHIMF"
*/
func TestGetBigipRegKey(t *testing.T) {

	url := "https://192.168.45.52//mgmt/tm/shared/licensing/registration"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Creating new HTTP request error: %v ", err)
	}

	t.Logf("Posting GET BIGIP Reg Key request on %v", url)
	req.SetBasicAuth("admin", "admin")
	httpClient = createHTTPClient()
	httpResp, responseMap := httpReq(req)
	if httpResp == nil || responseMap == nil {
		t.Errorf("HTTP response nil")
	}
	if responseMap["registrationKey"] != nil {
		registrationKey := responseMap["registrationKey"].(string)
		t.Logf("registrationKey: %s", registrationKey)
	}
}

type AS3Config struct {
	resourceConfig        as3ADC
	configmaps            []*AS3ConfigMap
	overrideConfigmapData string
	tenantMap             map[string]interface{}
	unifiedDeclaration    as3Declaration
}

type AS3ConfigMap struct {
	Name      string   // AS3 specific ConfigMap name
	Namespace string   // AS3 specific ConfigMap namespace
	config    as3ADC   // if AS3 Name is present, populate this with AS3 template data.
	endPoints []Member // Endpoints of all the pools in the configmap
	Validated bool     // Json Schema validated ok
}

func generateAS3ResourceDeclaration() as3ADC {

	app := as3Application{}
	app["class"] = "Application"
	app["template"] = "shared"

	tnt := as3Tenant{}
	tnt["class"] = "Tenant"
	tnt["Shared"] = app
	tnt["defaultRouteDomain"] = 0

	adc := as3ADC{}
	adc["k8s"] = tnt

	return adc
}

func prepareAS3ResourceConfig() as3ADC {

	adc := generateAS3ResourceDeclaration()
	controlObj := make(as3Control)
	controlObj["class"] = "Controls"
	controlObj["userAgent"] = "CIS/v K8S/v 1.23.10"
	adc["controls"] = controlObj

	return adc
}

func assertToBe(kind string, obj interface{}) bool {
	if obj == nil {
		return false
	}
	return (reflect.TypeOf(obj).Kind().String() == kind)
}

func getClass(obj interface{}) string {
	cfg, ok := obj.(map[string]interface{})
	if !ok {
		// If not a json object it doesn't have class attribute
		return ""
	}
	cl, ok := cfg["class"]
	if !ok {
		fmt.Println("No class attribute found")
		return ""
	}
	return cl.(string)
}

func getAS3ObjectFromTemplate(template as3Template) (as3Object, bool) {

	var tmpl map[string]interface{}
	err := json.Unmarshal([]byte(template), &tmpl)
	if err != nil {
		fmt.Errorf("[AS3] JSON unmarshal failed: %v  %v", err, template)
		return nil, false
	}

	as3 := make(as3Object)
	dclr := tmpl["declaration"]
	if dclr == nil || !assertToBe("map", dclr) {
		fmt.Println("[AS3] No ADC class declaration found or with wrong content.")
		return nil, false
	}

	for tn, t := range dclr.(map[string]interface{}) {
		if !assertToBe("map", t) {
			continue
		}
		tnt := t.(map[string]interface{})
		if tnt["class"] != "Tenant" {
			continue
		}
		as3[tenantName(tn)] = make(tenant, 0)
		for an, a := range t.(map[string]interface{}) {
			if !assertToBe("map", a) {
				continue
			}
			as3[tenantName(tn)][appName(an)] = []poolName{}
			for pn, v := range a.(map[string]interface{}) {
				if !assertToBe("map", v) {
					continue
				}
				if cl := getClass(v); cl != "Pool" {
					continue
				}
				mems := (v.(map[string]interface{}))["members"]
				if mems == nil {
					continue
				}
				if !assertToBe("slice", mems) || len(mems.([]interface{})) == 0 {
					continue
				}
				if !assertToBe("map", (mems.([]interface{}))[0]) {
					continue
				}
				mem0 := (mems.([]interface{}))[0].(map[string]interface{})
				srvAddrs := mem0["serverAddresses"]
				if srvAddrs == nil || len(srvAddrs.([]interface{})) != 0 {
					continue
				}
				as3[tenantName(tn)][appName(an)] = append(
					as3[tenantName(tn)][appName(an)],
					poolName(pn),
				)
			}
			if len(as3[tenantName(tn)][appName(an)]) == 0 {
				fmt.Printf("[AS3] No pools declared for application: %s, tenant: %s\n", an, tn)
			}
		}
	}
	if len(as3) == 0 {
		fmt.Println("[AS3] No tenants declared in AS3 template")
		return as3, false
	}
	return as3, true
}

func GetAppEndpoints() []Member {

	var members []Member

	member1 := Member{
		Address: generateIP(),
		Port:    8080,
		SvcPort: 80,
	}
	member2 := Member{
		Address: generateIP(),
		Port:    8080,
		SvcPort: 80,
	}

	member3 := Member{
		Address: generateIP(),
		Port:    8080,
		SvcPort: 80,
	}

	member4 := Member{
		Address: generateIP(),
		Port:    8080,
		SvcPort: 80,
	}

	member5 := Member{
		Address: generateIP(),
		Port:    8080,
		SvcPort: 80,
	}

	members = append(members, member1)
	members = append(members, member2)
	members = append(members, member3)
	members = append(members, member4)
	members = append(members, member5)

	return members
}

func processCfgMap(name string) (map[string]interface{}, []Member) {

	data, err := test.LoadFileAsString(name)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil
	}
	as3Tmpl := as3Template(data)
	obj, ok := getAS3ObjectFromTemplate(as3Tmpl)
	if !ok {
		fmt.Println("[AS3][Configmap] Error processing AS3 template")
		fmt.Printf("[AS3]Error in processing the ConfigMap: %v/%v", "f5-hub-1", "cm-cistest")
		return nil, nil
	}

	if _, ok := obj[tenantName(DEFAULT_PARTITION)]; ok {
		fmt.Printf("[AS3] Error in processing the ConfigMap: %v/%v", "f5-hub-1", "cm-cistest")
		fmt.Printf("[AS3] CIS managed partition <%s> should not be used in ConfigMaps as a Tenant", DEFAULT_PARTITION)
		return nil, nil
	}

	var tmp interface{}

	err = json.Unmarshal([]byte(as3Tmpl), &tmp)
	if nil != err {
		return nil, nil
	}

	templateJSON := tmp.(map[string]interface{})
	dec := (templateJSON["declaration"]).(map[string]interface{})
	tenantMap := make(map[string]interface{})
	var members []Member

	for tnt, apps := range obj {
		tenantObj := dec[string(tnt)].(map[string]interface{})
		tenantObj["defaultRouteDomain"] = 0
		for app, pools := range apps {
			appObj := tenantObj[string(app)].(map[string]interface{})
			for _, pn := range pools {
				poolObj := appObj[string(pn)].(map[string]interface{})
				eps := GetAppEndpoints()
				if len(eps) == 0 {
					continue
				}
				poolMem := (((poolObj["members"]).([]interface{}))[0]).(map[string]interface{})
				var ips []string
				var port int32
				for _, v := range eps {
					if int(v.SvcPort) == int(poolMem["servicePort"].(float64)) {
						ips = append(ips, v.Address)
						members = append(members, v)
						port = v.Port
					}
				}
				if port == 0 {
					ipMap := make(map[string]bool)
					members = append(members, eps...)
					for _, v := range eps {
						if _, ok := ipMap[v.Address]; !ok {
							ipMap[v.Address] = true
							ips = append(ips, v.Address)
						}
					}
					port = eps[0].Port
				}
				poolMem["serverAddresses"] = ips
				poolMem["servicePort"] = port
			}
		}
		tenantMap[string(tnt)] = tenantObj
	}
	return tenantMap, members
}

func prepareResourceAS3ConfigMaps(name string) ([]*AS3ConfigMap, string) {

	var as3Cfgmaps []*AS3ConfigMap
	var overriderAS3CfgmapData string

	cfgmap := &AS3ConfigMap{
		Name:      "cm-cistest",
		Namespace: "f5-hub-1",
		Validated: true,
	}

	tenantMap, endPoints := processCfgMap(name)
	cfgmap.config = tenantMap
	cfgmap.endPoints = endPoints
	as3Cfgmaps = append(as3Cfgmaps, cfgmap)

	return as3Cfgmaps, overriderAS3CfgmapData
}

func updateTenantMap(tempAS3Config AS3Config) AS3Config {
	// Parse as3Config.configmaps , extract all tenants and store in tenantMap.
	for _, cm := range tempAS3Config.configmaps {
		for tenantName, tenant := range cm.config {
			tempAS3Config.tenantMap[tenantName] = tenant
		}
	}
	return tempAS3Config
}

func getADC() map[string]interface{} {
	var as3Obj map[string]interface{}

	baseAS3ConfigTemplate := fmt.Sprintf(baseAS3Config, "3.36.1", "3.36.1-1", "3.36.0")
	_ = json.Unmarshal([]byte(baseAS3ConfigTemplate), &as3Obj)

	return as3Obj
}

func prepareTenantDeclaration(cfg *AS3Config, tenantName string) as3Declaration {

	as3Obj := getADC()
	adc, _ := as3Obj["declaration"].(map[string]interface{})

	adc[tenantName] = cfg.tenantMap[tenantName]

	unifiedDecl, err := json.Marshal(as3Obj)
	if err != nil {
		fmt.Printf("[AS3] Unified declaration: %v", err)
	}

	return as3Declaration(unifiedDecl)
}

func getAS3APIURL(tenants []string) string {
	apiURL := "https://192.168.45.52/mgmt/shared/appsvcs/declare/" + strings.Join(tenants, ",")
	return apiURL
}

/*
The steps for create first VS in a clean BIG-IP VE:

 1. Prepare as3 Configmap base on appMgr's RequestMessage via channel, the sub tasks of prepare as3 Configmap including:

 1. extract the service endpoint

 2. Preapre Declaration

 3. Create HTTP Client

 4. Process HTTP POST request and verify the results

    This test func can be simulated via curl:

    % curl -s -k -u "admin:admin" -X POST -H "Content-Type: application/json" -d @$(pwd)/cm1.json https://192.168.45.52/mgmt/shared/appsvcs/declare/cistest001

    {
    "results": [
    {
    "code": 200,
    "message": "success",
    "lineCount": 25,
    "host": "localhost",
    "tenant": "cistest001",
    "runTime": 2189
    }
    ],
    "declaration": {
    "cistest001": {
    "app-1": {
    "app_1_svc_pool": {
    "class": "Pool",
    "loadBalancingMode": "least-connections-member",
    "members": [
    {
    "serverAddresses": [
    "1.1.1.1",
    "1.1.1.2"
    ],
    "servicePort": 8080
    }
    ],
    "monitors": [
    "tcp"
    ]
    },
    "app_svc_vs": {
    "class": "Service_HTTP",
    "persistenceMethods": [
    "cookie"
    ],
    "pool": "app_1_svc_pool",
    "snat": "self",
    "virtualAddresses": [
    "192.168.200.32"
    ],
    "virtualPort": 80
    },
    "class": "Application",
    "template": "generic"
    },
    "class": "Tenant",
    "defaultRouteDomain": 0
    },
    "class": "ADC",
    "controls": {
    "class": "Controls",
    "userAgent": "CIS Configured AS3",
    "archiveTimestamp": "2024-01-14T07:47:11.141Z"
    },
    "id": "urn:uuid:85626792-9ee7-46bb-8fc8-4ba708cfdc1d",
    "label": "CIS Declaration",
    "remark": "Auto-generated by CIS",
    "schemaVersion": "3.36.0",
    "updateMode": "selective"
    }
    }
*/
func TestCreateFirstVSInCleanVEViaCM(t *testing.T) {

	as3Config := &AS3Config{
		tenantMap: make(map[string]interface{}),
	}

	// Process Route or Ingress
	as3Config.resourceConfig = prepareAS3ResourceConfig()

	// Process all Configmaps (including overrideAS3)
	as3Config.configmaps, as3Config.overrideConfigmapData = prepareResourceAS3ConfigMaps("cm1.txt")

	updateTenantMap(*as3Config)

	partition := "cistest001"
	tenantDecl := prepareTenantDeclaration(as3Config, partition)
	data := string(tenantDecl)

	t.Logf("declaration: %s", data)
	url := getAS3APIURL([]string{partition})
	t.Logf("url: %s", url)

	cfg := config{
		data:      data,
		as3APIURL: url,
	}
	httpReqBody := bytes.NewBuffer([]byte(cfg.data))

	req, err := http.NewRequest("POST", cfg.as3APIURL, httpReqBody)
	if err != nil {
		t.Errorf("[AS3] Creating new HTTP request error: %v ", err)
	}
	t.Logf("[AS3] posting request to %v", cfg.as3APIURL)

	req.SetBasicAuth("admin", "admin")
	httpClient = createHTTPClient()
	httpResp, responseMap := httpReq(req)
	t.Logf("httpResp: %v", httpResp)
	t.Logf("response: %v", responseMap)

	results := (responseMap["results"]).([]interface{})
	for _, value := range results {
		v := value.(map[string]interface{})
		t.Logf("[AS3] Response from BIG-IP: code: %v, tenant: %v, message: %v, runTime: %v", v["code"], v["tenant"], v["message"], v["runTime"])
	}

}

func getEmptyAs3Declaration(partition string) as3Declaration {

	var as3Config map[string]interface{}
	baseAS3ConfigEmpty := fmt.Sprintf(baseAS3Config, "3.36.1", "3.36.1-1", "3.36.0")
	_ = json.Unmarshal([]byte(baseAS3ConfigEmpty), &as3Config)
	decl := as3Config["declaration"].(map[string]interface{})

	controlObj := make(as3Control)
	controlObj["class"] = "Controls"
	controlObj["userAgent"] = "CIS/v K8S/v 1.23.10"
	decl["controls"] = controlObj

	decl[partition] = map[string]string{"class": "Tenant"}
	data, _ := json.Marshal(as3Config)
	emptyAS3Declaration := as3Declaration(data)
	return emptyAS3Declaration
}

/*
Steps to delete a partition

 1. Create a empty declaration file
 2. Create a HTTP Client
 3. Execute HTTP POST with empty declaration file as body
 4. Check the HTTP Response

This test function can be simulated with the below curl:

% curl -s -k -u "admin:admin" -X POST -H "Content-Type: application/json" -d @$(pwd)/cm1-del.json https://192.168.45.52/mgmt/shared/appsvcs/declare/cistest001

	{
	  "results": [
	    {
	      "code": 200,
	      "message": "success",
	      "lineCount": 30,
	      "host": "localhost",
	      "tenant": "cistest001",
	      "runTime": 1764
	    }
	  ],
	  "declaration": {
	    "class": "ADC",
	    "controls": {
	      "class": "Controls",
	      "userAgent": "CIS/v K8S/v1.23.10",
	      "archiveTimestamp": "2024-01-14T08:18:45.228Z"
	    },
	    "id": "urn:uuid:85626792-9ee7-46bb-8fc8-4ba708cfdc1d",
	    "label": "CIS Declaration",
	    "remark": "Auto-generated by CIS",
	    "schemaVersion": "3.36.0",
	    "updateMode": "selective"
	  }
	}
*/
func TestDeletePartition(t *testing.T) {

	partition := "cistest001"
	process_deletion(partition)

}

func TestCreateMultipleVS(t *testing.T) {

	globalMap = make(map[string]bool)

	as3Config := &AS3Config{
		tenantMap: make(map[string]interface{}),
	}

	// Process Route or Ingress
	as3Config.resourceConfig = prepareAS3ResourceConfig()

	// Process all Configmaps (including overrideAS3)
	as3Config.configmaps, as3Config.overrideConfigmapData = prepareResourceAS3ConfigMaps("cm2.txt")

	updateTenantMap(*as3Config)

	t.Logf("as3Config: %v", as3Config)

	for partition, tenant := range as3Config.tenantMap {
		t.Logf("%s, %v", partition, tenant)
		tenantDecl := prepareTenantDeclaration(as3Config, partition)
		data := string(tenantDecl)
		url := getAS3APIURL([]string{partition})
		t.Logf("data: %s", data)
		t.Logf("url: %s", url)
		cfg := config{
			data:      data,
			as3APIURL: url,
		}
		post_as3_declaration(cfg)
	}

}

func TestDeleteMultipleVS(t *testing.T) {
	globalMap = make(map[string]bool)
	partitions := exreact_tenant_names("cm2.txt")
	for _, partition := range partitions {
		process_deletion(partition)
	}
}
