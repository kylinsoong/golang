package channel

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/kylinsoong/bigip-ctlr/test"
)

type ResourceMap map[int32][]*ResourceConfig

type MessageRequest struct {
	ReqID   uint
	MsgType string
	ResourceRequest
}

type ResourceRequest struct {
	PoolMembers  map[Member]struct{}
	Resources    *AgentResources
	Profs        map[SecretKey]CustomProfile
	IrulesMap    IRulesMap
	IntDgMap     InternalDataGroupMap
	IntF5Res     InternalF5ResourcesGroup
	AgentCfgmaps []*AgentCfgMap
}

type AgentCfgMap struct {
	Operation    string
	GetEndpoints func(string, string) []Member
	Data         string
	Name         string
	Namespace    string
	Label        map[string]string
}

type InternalF5ResourcesGroup map[string]InternalF5Resources
type InternalF5Resources map[Record]F5Resources

type F5Resources struct {
	Virtual   ConstVirtuals // 0 - HTTP, 1 - HTTPS, 2 - HTTP/S
	WAFPolicy string
}

type ConstVirtuals int

type Record struct {
	Host string
	Path string
}

type InternalDataGroupMap map[NameRef]DataGroupNamespaceMap

type DataGroupNamespaceMap map[string]*InternalDataGroup

type InternalDataGroup struct {
	Name      string                   `json:"name"`
	Partition string                   `json:"-"`
	Records   InternalDataGroupRecords `json:"records"`
}

type InternalDataGroupRecords []InternalDataGroupRecord

type InternalDataGroupRecord struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type IRulesMap map[NameRef]*IRule

type IRule struct {
	Name      string `json:"name"`
	Partition string `json:"-"`
	Code      string `json:"apiAnonymous"`
}

type SecretKey struct {
	Name         string
	ResourceName string
}

type CustomProfile struct {
	Name         string `json:"name"`
	Partition    string `json:"-"`
	Context      string `json:"context"` // 'clientside', 'serverside', or 'all'
	Cert         string `json:"cert"`
	Key          string `json:"key"`
	ServerName   string `json:"serverName,omitempty"`
	SNIDefault   bool   `json:"sniDefault,omitempty"`
	PeerCertMode string `json:"peerCertMode,omitempty"`
	CAFile       string `json:"caFile,omitempty"`
	ChainCA      string `json:"chainCA,onitempty"`
}

type AgentResources struct {
	RsMap  ResourceConfigMap
	RsCfgs ResourceConfigs
}

type ResourceConfigs []*ResourceConfig

type ResourceConfigMap map[string]*ResourceConfig

type ResourceConfig struct {
	MetaData       MetaData         `json:"-"`
	Virtual        Virtual          `json:"virtual,omitempty"`
	IApp           IApp             `json:"iapp,omitempty"`
	Pools          Pools            `json:"pools,omitempty"`
	Monitors       Monitors         `json:"monitors,omitempty"`
	Policies       Policies         `json:"policies,omitempty"`
	ServiceAddress []ServiceAddress `json:"serviceAddress,omitempty"`
}

type ServiceAddress struct {
	ArpEnabled         bool   `json:"arpEnabled,omitempty"`
	ICMPEcho           string `json:"icmpEcho,omitempty"`
	RouteAdvertisement string `json:"routeAdvertisement,omitempty"`
	TrafficGroup       string `json:"trafficGroup,omitempty,omitempty"`
	SpanningEnabled    bool   `json:"spanningEnabled,omitempty"`
}

type Policies []Policy

type Policy struct {
	Name        string   `json:"name"`
	Partition   string   `json:"-"`
	SubPath     string   `json:"subPath,omitempty"`
	Controls    []string `json:"controls,omitempty"`
	Description string   `json:"description,omitempty"`
	Legacy      bool     `json:"legacy,omitempty"`
	Requires    []string `json:"requires,omitempty"`
	Rules       Rules    `json:"rules,omitempty"`
	Strategy    string   `json:"strategy,omitempty"`
}

type Rules []*Rule

type Rule struct {
	Name       string       `json:"name"`
	FullURI    string       `json:"-"`
	Ordinal    int          `json:"ordinal,omitempty"`
	Actions    []*Action    `json:"actions,omitempty"`
	Conditions []*Condition `json:"conditions,omitempty"`
}

type Action struct {
	Name      string `json:"name"`
	Pool      string `json:"pool,omitempty"`
	HTTPHost  bool   `json:"httpHost,omitempty"`
	HttpReply bool   `json:"httpReply,omitempty"`
	HTTPURI   bool   `json:"httpUri,omitempty"`
	Forward   bool   `json:"forward,omitempty"`
	Location  string `json:"location,omitempty"`
	Path      string `json:"path,omitempty"`
	Redirect  bool   `json:"redirect,omitempty"`
	Replace   bool   `json:"replace,omitempty"`
	Request   bool   `json:"request,omitempty"`
	Reset     bool   `json:"reset,omitempty"`
	Select    bool   `json:"select,omitempty"`
	Value     string `json:"value,omitempty"`
}

type Condition struct {
	Name            string   `json:"name"`
	Address         bool     `json:"address,omitempty"`
	CaseInsensitive bool     `json:"caseInsensitive,omitempty"`
	Equals          bool     `json:"equals,omitempty"`
	EndsWith        bool     `json:"endsWith,omitempty"`
	External        bool     `json:"external,omitempty"`
	HTTPHost        bool     `json:"httpHost,omitempty"`
	Host            bool     `json:"host,omitempty"`
	HTTPURI         bool     `json:"httpUri,omitempty"`
	Index           int      `json:"index,omitempty"`
	Matches         bool     `json:"matches,omitempty"`
	Path            bool     `json:"path,omitempty"`
	PathSegment     bool     `json:"pathSegment,omitempty"`
	Present         bool     `json:"present,omitempty"`
	Remote          bool     `json:"remote,omitempty"`
	Request         bool     `json:"request,omitempty"`
	Scheme          bool     `json:"scheme,omitempty"`
	Tcp             bool     `json:"tcp,omitempty"`
	Values          []string `json:"values"`
}

type Monitors []Monitor

type Monitor struct {
	Name      string `json:"name"`
	Partition string `json:"-"`
	Interval  int    `json:"interval,omitempty"`
	Type      string `json:"type,omitempty"`
	Send      string `json:"send,omitempty"`
	Recv      string `json:"recv,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
}

type Pools []Pool

type Pool struct {
	Name         string   `json:"name"`
	Partition    string   `json:"-"`
	ServiceName  string   `json:"-"`
	ServicePort  int32    `json:"-"`
	Balance      string   `json:"loadBalancingMode"`
	Members      []Member `json:"members"`
	MonitorNames []string `json:"monitors,omitempty"`
}

type IApp struct {
	Name                string                    `json:"name"`
	Partition           string                    `json:"-"`
	IApp                string                    `json:"template"`
	IAppPoolMemberTable *iappPoolMemberTable      `json:"poolMemberTable,omitempty"`
	IAppOptions         map[string]string         `json:"options,omitempty"`
	IAppTables          map[string]iappTableEntry `json:"tables,omitempty"`
	IAppVariables       map[string]string         `json:"variables,omitempty"`
}

type iappTableEntry struct {
	Columns []string   `json:"columns,omitempty"`
	Rows    [][]string `json:"rows,omitempty"`
}

type iappPoolMemberTable struct {
	Name    string                 `json:"name"`
	Columns []iappPoolMemberColumn `json:"columns"`
	Members []Member               `json:"members,omitempty"`
}

type iappPoolMemberColumn struct {
	Name  string `json:"name"`
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
}

type Virtual struct {
	Name                  string                `json:"name"`
	PoolName              string                `json:"pool,omitempty"`
	Partition             string                `json:"-"`
	Destination           string                `json:"destination"`
	Enabled               bool                  `json:"enabled"`
	IpProtocol            string                `json:"ipProtocol,omitempty"`
	SourceAddrTranslation SourceAddrTranslation `json:"sourceAddressTranslation,omitempty"`
	Policies              []NameRef             `json:"policies,omitempty"`
	IRules                []string              `json:"rules,omitempty"`
	Profiles              ProfileRefs           `json:"profiles,omitempty"`
	Description           string                `json:"description,omitempty"`
	VirtualAddress        *VirtualAddress       `json:"-"`
}

type SourceAddrTranslation struct {
	Type string `json:"type"`
	Pool string `json:"pool,omitempty"`
}

type NameRef struct {
	Name      string `json:"name"`
	Partition string `json:"partition"`
}

type VirtualAddress struct {
	BindAddr string `json:"bindAddr,omitempty"`
	Port     int32  `json:"port,omitempty"`
}

type ProfileRefs []ProfileRef

type ProfileRef struct {
	Name      string `json:"name"`
	Partition string `json:"partition"`
	Context   string `json:"context"`
	Namespace string `json:"-"`
}

type MetaData struct {
	Active       bool
	ResourceType string
	RouteProfs   map[RouteKey]string
	IngName      string
}

type RouteKey struct {
	Name      string
	Namespace string
	Context   string
}

type Member struct {
	Address string `json:"address"`
	Port    int32  `json:"port"`
	SvcPort int32  `json:"svcPort"`
	Session string `json:"session,omitempty"`
}

func numberGenerator(inputCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 5; i++ {
		fmt.Printf("numberGenerator write %d to %v\n", i, inputCh)
		inputCh <- i // Send numbers 1 to 5 to the channel
	}
	close(inputCh) // Close the channel to signal no more data will be sent
}

func squareCalculator(ch chan int, resultCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("squareCalculator block on %v\n", ch)
	for num := range ch {
		square := num * num
		fmt.Printf("squareCalculator write %d to %v\n", square, resultCh)
		resultCh <- square // Send squared result to the resultCh channel
	}
	close(resultCh) // Close the resultCh channel to signal no more results will be sent
}

func resultPrinter(resultCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("resultPrinter block on %v\n", resultCh)
	for result := range resultCh {
		fmt.Println("Squared Result:", result)
	}
}

func worker(ch chan struct{}) {
	fmt.Printf("Worker is starting...(%v)\n", ch)
	time.Sleep(2 * time.Second)
	fmt.Println("Worker is done!")
	ch <- struct{}{}
}

func Test_Communication_Synchronization(t *testing.T) {

	numberCh := make(chan int)
	resultCh := make(chan int)
	fmt.Printf("Main Thread Start, created %v, %v\n", numberCh, resultCh)
	var wg sync.WaitGroup
	wg.Add(3)
	go resultPrinter(resultCh, &wg)
	time.Sleep(1000 * time.Millisecond)
	go squareCalculator(numberCh, resultCh, &wg)
	time.Sleep(1000 * time.Millisecond)
	go numberGenerator(numberCh, &wg)
	time.Sleep(1000 * time.Millisecond)
	wg.Wait()
}

func prepareResourceKey(kind, namespace, name string) string {
	return kind + "_" + namespace + "/" + name
}

func validateConfigJson(tmpConfig string) error {
	var tmp interface{}
	err := json.Unmarshal([]byte(tmpConfig), &tmp)
	return err
}

func getEndpoints(selector, namespace string) []Member {
	return nil
}

func GetAllResources() ResourceConfigs {
	var cfgs ResourceConfigs
	return cfgs
}

func createMessageRequest(filePath string) MessageRequest {

	resKey := prepareResourceKey("configmaps", "f5-hub-1", "cm-cistest")

	svcPortMap := make(map[int32]bool)
	rsMap := make(ResourceMap)
	dgMap := make(InternalDataGroupMap)

	fmt.Printf("%v, svcPortMap: %v, rsMap: %v, dgMap: %v\n", resKey, svcPortMap, rsMap, dgMap)

	fileContent, err := test.ReadFileAsString(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	validateConfigJson(fileContent)

	labels := map[string]string{
		"as3":    "true",
		"f5type": "virtual-server",
	}
	agntCfgMap := new(AgentCfgMap)
	agntCfgMap.Name = "cm-cistest"
	agntCfgMap.Namespace = "f5-hub-1"
	agntCfgMap.Data = fileContent
	agntCfgMap.Label = labels
	agntCfgMap.GetEndpoints = getEndpoints

	Profs := map[SecretKey]CustomProfile{}
	idgMap := make(InternalDataGroupMap)
	iRulesMap := make(IRulesMap)
	resourceConfigMap := make(ResourceConfigMap)
	resourceConfigs := GetAllResources()
	intF5Res := InternalF5ResourcesGroup{}

	agentCfgMapLst := []*AgentCfgMap{}
	agentCfgMapLst = append(agentCfgMapLst, agntCfgMap)

	deployCfg := ResourceRequest{
		Resources: &AgentResources{
			RsMap:  resourceConfigMap,
			RsCfgs: resourceConfigs,
		},
		Profs:        Profs,
		IrulesMap:    iRulesMap,
		IntDgMap:     idgMap,
		IntF5Res:     intF5Res,
		AgentCfgmaps: agentCfgMapLst,
	}

	agentReq := MessageRequest{MsgType: "L4L7Declaration", ResourceRequest: deployCfg}

	return agentReq
}

func Deploy(reqChan chan MessageRequest, req interface{}) error {
	msgReq := req.(MessageRequest)
	select {
	case reqChan <- msgReq:
	case <-reqChan:
		reqChan <- msgReq
	}
	fmt.Printf("[AS3] Sent message to %v, ReqID: %d, MsgType: %s, ResourceRequest: %v\n", reqChan, msgReq.ReqID, msgReq.MsgType, msgReq.ResourceRequest)
	return nil
}

func Test_Signal_Completion(t *testing.T) {
	doneCh := make(chan struct{})
	fmt.Printf("Main function starting...(%v)\n", doneCh)
	go worker(doneCh)
	<-doneCh
	fmt.Println("Main function exiting.")
}

func Test_Signal_Interruptiong_Termination(t *testing.T) {
	fmt.Println("Started to run tasks...")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	fmt.Printf("Received signal: %v\n", sig)
}

func Test_Send_MessageRequest_via_Chan(t *testing.T) {

	filePath := "/Users/k.song/src/golang/bigip-ctlr/test/configmaps/cm3.txt"
	agentReq := createMessageRequest(filePath)
	fmt.Printf("%v\n", len(agentReq.AgentCfgmaps))
	ReqChan := make(chan MessageRequest, 1)
	go func() {
		for req := range ReqChan {
			for _, agentCfgMap := range req.AgentCfgmaps {
				fmt.Printf("name: %s, namespace: %s, data:\n", agentCfgMap.Name, agentCfgMap.Namespace)
				fmt.Println(agentCfgMap.Data)
			}

		}
	}()
	err := Deploy(ReqChan, agentReq)
	if err != nil {
		fmt.Println("Error deploying request:", err)
	}
}
