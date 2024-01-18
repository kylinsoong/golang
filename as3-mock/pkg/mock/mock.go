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

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    log.Infof("Received POST request from %s for %s", r.RemoteAddr, r.URL.Path)
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
    log.Infof("Request Body: \n%+v", prettyJSON.String())

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

    registrationInfo := RegistrationInfo{
	Vendor:               "F5 Networks, Inc.",
	LicensedDateTime:     "2024-01-10T00:00:00-08:00",
	LicensedVersion:      "15.1.10",
	LicenseEndDateTime:   "2025-02-10T00:00:00-08:00",
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

func HandleGetInfo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

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
