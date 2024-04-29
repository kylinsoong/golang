package mock

import (
    "encoding/json"
    "net/http"
    "io/ioutil"
    "bytes"
    "strings"

    log  "github.com/kylinsoong/golang/as3-mock/pkg/vlogger"
)

type InfoResponse struct {
    Version        string `json:"version"`
    Release        string `json:"release"`
    SchemaCurrent  string `json:"schemaCurrent"`
    SchemaMinimum  string `json:"schemaMinimum"`
}

type RegistrationInfo struct {
    Vendor               string `json:"vendor"`
    LicensedDateTime     string `json:"licensedDateTime"`
    LicensedVersion      string `json:"licensedVersion"`
    LicenseEndDateTime   string `json:"licenseEndDateTime"`
    LicenseStartDateTime string `json:"licenseStartDateTime"`
    RegistrationKey      string `json:"registrationKey"`
}

type Token struct {
    Token            string `json:"token"`
    Name             string `json:"name"`
    UserName         string `json:"userName"`
    AuthProviderName string `json:"authProviderName"`
    User             struct {
        Link string `json:"link"`
    } `json:"user"`
    Timeout          int    `json:"timeout"`
    StartTime        string `json:"startTime"`
    Address          string `json:"address"`
    Partition        string `json:"partition"`
    Generation       int    `json:"generation"`
    LastUpdateMicros int64  `json:"lastUpdateMicros"`
    ExpirationMicros int64  `json:"expirationMicros"`
    Kind             string `json:"kind"`
    SelfLink         string `json:"selfLink"`
}

type AuthResponse struct {
    Username          string  `json:"username"`
    LoginReference    struct {
        Link string `json:"link"`
    } `json:"loginReference"`
    LoginProviderName string `json:"loginProviderName"`
    Token             Token  `json:"token"`
    Generation        int    `json:"generation"`
    LastUpdateMicros  int64  `json:"lastUpdateMicros"`
}

type Reference struct {
    Link string `json:"link"`
}

type Item struct {
    Reference Reference `json:"reference"`
}

type SysCollectionState struct {
    Kind     string `json:"kind"`
    SelfLink string `json:"selfLink"`
    Items    []Item `json:"items"`
}


func HandlePostRequestAuthToken(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var reqBody struct {
        Username          string `json:"username"`
        Password          string `json:"password"`
        LoginProviderName string `json:"loginProviderName"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    log.Infof("Received POST request from %s for %s, Request Body: username %s, password: %s, loginProviderName: %s", r.RemoteAddr, r.URL.Path, reqBody.Username, reqBody.Password, reqBody.LoginProviderName)

    authResp := AuthResponse{
        Username:          reqBody.Username,
        LoginProviderName: "tmos",
        Token: Token{
            Token:            "QLLJFAT5DHOGEGE2XCAX6ZBEKA",
            Name:             "QLLJFAT5DHOGEGE2XCAX6ZBEKA",
            UserName:         reqBody.Username,
            AuthProviderName: "tmos",
            User: struct {
                Link string `json:"link"`
            }{
                Link: "https://localhost/mgmt/shared/authz/users/admin",
            },
            Timeout:          1200,
            StartTime:        "2024-04-29T19:17:15.460+0800",
            Address:          "10.1.10.240",
            Partition:        "[All]",
            Generation:       1,
            LastUpdateMicros: 1714389435460230,
            ExpirationMicros: 1814390635460000,
            Kind:             "shared:authz:tokens:authtokenitemstate",
            SelfLink:         "https://localhost/mgmt/shared/authz/tokens/QLLJFAT5DHOGEGE2XCAX6ZBEKA",
        },
        Generation:       0,
        LastUpdateMicros: 0,
    }

    respBody, err := json.Marshal(authResp)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(respBody)
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusInternalServerError)
	log.Errorf("Error reading request body: %v", err)
	return
    }
    defer r.Body.Close()

    var prettyJSON bytes.Buffer
    err = json.Indent(&prettyJSON, body, "", "    ")
    if err != nil {
        log.Errorf("Error formatting JSON: %v", err)
        return
    }
    log.Infof("Received POST request from %s for %s, Request Body: \n%+v", r.RemoteAddr, r.URL.Path, prettyJSON.String())

    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) != 6 || pathParts[4] != "declare" {
	http.Error(w, "Invalid path", http.StatusBadRequest)
	return
    }

    tenant := pathParts[5]

    response := map[string]interface{}{
	"results": []map[string]interface{}{
	    {
		"code":      200,
		"message":   "success",
		"lineCount": 30,
		"host":      "localhost",
		"tenant":    tenant,
		"runTime":   100,
	    },
	},
	"declaration": map[string]interface{}{
	    "class":      "ADC",
	    "updateMode": "selective",
	},
    }

    responseJSON, err := json.Marshal(response)
    if err != nil {
	http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseJSON)
}

func HandleGetRegistration(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    log.Infof("Received GET request from %s for %s", r.RemoteAddr, r.URL.Path)

    registrationInfo := RegistrationInfo{
	Vendor:               "F5 Networks, Inc.",
	LicensedDateTime:     "2024-01-10T00:00:00-08:00",
	LicensedVersion:      "15.1.10",
	LicenseEndDateTime:   "2026-02-10T00:00:00-08:00",
	LicenseStartDateTime: "2024-01-09T00:00:00-08:00",
	RegistrationKey:      "KVPKO-EBYPF-UFQQG-WYBNP-TXRHIMF",
    }

    responseJSON, err := json.Marshal(registrationInfo)
    if err != nil {
	http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseJSON)
}

func HandleGetTMOSSys(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    log.Infof("Received GET request from %s for %s", r.RemoteAddr, r.URL.Path)

    items := []Item{
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/application?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/crypto?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/daemon-log-settings?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/diags?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/disk?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/dynad?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/ecm?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/file?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/fpga?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/icall?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/log-config?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/pfman?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/sflow?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/software?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/turboflex?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/url-db?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/aom?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/autoscale-group?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/cluster?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/config?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/core?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/daemon-ha?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/datastor?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/db?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/dns?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/feature-module?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/folder?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/global-settings?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/ha-group?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/httpd?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/icontrol-soap?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/internal-proxy?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/log-rotate?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/management-dhcp?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/management-ip?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/management-ovsdb?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/management-proxy-config?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/management-route?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/ntp?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/outbound-smtp?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/provision?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/scriptd?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/service?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/smtp-server?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/snmp?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/sshd?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/state-mirroring?ver=15.1.10"}},
	{Reference: Reference{Link: "https://localhost/mgmt/tm/sys/syslog?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/telemd?ver=15.1.10"}},
        {Reference: Reference{Link: "https://localhost/mgmt/tm/sys/ucs?ver=15.1.10"}},
    }

    sysCollectionState := SysCollectionState{
        Kind:     "tm:sys:syscollectionstate",
        SelfLink: "https://localhost/mgmt/tm/sys?ver=15.1.10",
        Items:    items,
    }

    jsonResponse, err := json.Marshal(sysCollectionState)
    if err != nil {
        http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}


func HandleGetInfo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    log.Infof("Received GET request from %s for %s", r.RemoteAddr, r.URL.Path)

    response := InfoResponse{
	Version:       "3.36.1",
	Release:       "1",
	SchemaCurrent: "3.36.0",
	SchemaMinimum: "3.0.0",
    }

    responseJSON, err := json.Marshal(response)
    if err != nil {
	http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseJSON)
}
