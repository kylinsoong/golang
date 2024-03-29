package mock

import (
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "io/ioutil"
    "bytes"
    "time"

    . "github.com/kylinsoong/golang/nim-mock/pkg/resources"

    log  "github.com/kylinsoong/golang/nim-mock/pkg/vlogger"
)


func GetInstanceGroupsSummary(w http.ResponseWriter, r *http.Request){

    log.Infof("Received %s request from %s for %s", r.Method, r.RemoteAddr, r.URL.Path)

    response := LoadInstanceGroups{
	Count: 2,
	Items: []InstanceGroupSimple{
	    {
		Description:   "",
		DisplayName:   "企业网银",
		Name:          "qywy",
		UID:           "c78208d2-f017-4a48-bab9-f2dc2432b8c9",
	    },
	    {
		Description:   "",
		DisplayName:   "微信银行",
		Name:          "wxyh",
		UID:           "6c81fd7d-fdc6-45f2-9c2b-666e0ca89897",
	    },
	},
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
	http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	return
    } 

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func GetInstanceGroupsDeployments(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    uid := vars["uid"]
    log.Infof("Received %s request from %s for %s, instance-group ID: %s", r.Method, r.RemoteAddr, r.URL.Path, uid)
    
    deploymentDetail := DeploymentDetail{
        Failure: []InstanceStatus{},
        Pending: []InstanceActivityStatus{},
        Success: []InstanceStatus{
            {Name: "server-c.com",},
            {Name: "server-d.com",},
        },
    }

    response := DeploymentDetails{
        CreateTime:    time.Now(),
        Details:       deploymentDetail,
        ID:            "f68f185c-04ea-46d7-9025-b293a06fe961",
        Message:       "Instance Group config successfully published to qywy",
        Status:        "successful",
        UpdateTime:    time.Now(),
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)    
}

func InstanceGroupsConfig(w http.ResponseWriter, r *http.Request){

    vars := mux.Vars(r)
    uid := vars["uid"]

    log.Infof("Received %s request from %s for %s, instance-group ID: %s", r.Method, r.RemoteAddr, r.URL.Path, uid)

    switch r.Method {
    case http.MethodGet:

        auxFiles := AuxData{
	    Files: []FileData{},
	    RootDir: "/",
	}

        configFiles := ConfigData {
            Files: []FileData{
                {
                    Contents: "CnVzZXIgIG5naW54Owp3b3JrZXJfcHJvY2Vzc2VzICBhdXRvOwoKZXJyb3JfbG9nICAvdmFyL2xvZy9uZ2lueC9lcnJvci5sb2cgbm90aWNlOwpwaWQgICAgICAgIC92YXIvcnVuL25naW54LnBpZDsKCgpldmVudHMgewogICAgd29ya2VyX2Nvbm5lY3Rpb25zICAxMDI0Owp9CgoKaHR0cCB7CiAgICBpbmNsdWRlICAgICAgIC9ldGMvbmdpbngvbWltZS50eXBlczsKICAgIGRlZmF1bHRfdHlwZSAgYXBwbGljYXRpb24vb2N0ZXQtc3RyZWFtOwoKICAgIGxvZ19mb3JtYXQgIG1haW4gICckcmVtb3RlX2FkZHIgLSAkcmVtb3RlX3VzZXIgWyR0aW1lX2xvY2FsXSAiJHJlcXVlc3QiICcKICAgICAgICAgICAgICAgICAgICAgICckc3RhdHVzICRib2R5X2J5dGVzX3NlbnQgIiRodHRwX3JlZmVyZXIiICcKICAgICAgICAgICAgICAgICAgICAgICciJGh0dHBfdXNlcl9hZ2VudCIgIiRodHRwX3hfZm9yd2FyZGVkX2ZvciInOwoKICAgIGFjY2Vzc19sb2cgIC92YXIvbG9nL25naW54L2FjY2Vzcy5sb2cgIG1haW47CgogICAgc2VuZGZpbGUgICAgICAgIG9uOwogICAgI3RjcF9ub3B1c2ggICAgIG9uOwoKICAgIGtlZXBhbGl2ZV90aW1lb3V0ICA2NTsKCiAgICAjZ3ppcCAgb247CgogICAgaW5jbHVkZSAvZXRjL25naW54L2NvbmYuZC8qLmNvbmY7Cn0KCgojIFRDUC9VRFAgcHJveHkgYW5kIGxvYWQgYmFsYW5jaW5nIGJsb2NrCiMKI3N0cmVhbSB7CiAgICAjIEV4YW1wbGUgY29uZmlndXJhdGlvbiBmb3IgVENQIGxvYWQgYmFsYW5jaW5nCgogICAgI3Vwc3RyZWFtIHN0cmVhbV9iYWNrZW5kIHsKICAgICMgICAgem9uZSB0Y3Bfc2VydmVycyA2NGs7CiAgICAjICAgIHNlcnZlciBiYWNrZW5kMS5leGFtcGxlLmNvbToxMjM0NTsKICAgICMgICAgc2VydmVyIGJhY2tlbmQyLmV4YW1wbGUuY29tOjEyMzQ1OwogICAgI30KCiAgICAjc2VydmVyIHsKICAgICMgICAgbGlzdGVuIDEyMzQ1OwogICAgIyAgICBzdGF0dXNfem9uZSB0Y3Bfc2VydmVyOwogICAgIyAgICBwcm94eV9wYXNzIHN0cmVhbV9iYWNrZW5kOwogICAgI30KI30K",
                    Name: "/etc/nginx/nginx.conf",
                },
                {
                    Contents: "CnR5cGVzIHsKICAgIHRleHQvaHRtbCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBodG1sIGh0bSBzaHRtbDsKICAgIHRleHQvY3NzICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBjc3M7CiAgICB0ZXh0L3htbCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgeG1sOwogICAgaW1hZ2UvZ2lmICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIGdpZjsKICAgIGltYWdlL2pwZWcgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBqcGVnIGpwZzsKICAgIGFwcGxpY2F0aW9uL2phdmFzY3JpcHQgICAgICAgICAgICAgICAgICAgICAgICAgICBqczsKICAgIGFwcGxpY2F0aW9uL2F0b20reG1sICAgICAgICAgICAgICAgICAgICAgICAgICAgICBhdG9tOwogICAgYXBwbGljYXRpb24vcnNzK3htbCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJzczsKCiAgICB0ZXh0L21hdGhtbCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgbW1sOwogICAgdGV4dC9wbGFpbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHR4dDsKICAgIHRleHQvdm5kLnN1bi5qMm1lLmFwcC1kZXNjcmlwdG9yICAgICAgICAgICAgICAgICBqYWQ7CiAgICB0ZXh0L3ZuZC53YXAud21sICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgd21sOwogICAgdGV4dC94LWNvbXBvbmVudCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIGh0YzsKCiAgICBpbWFnZS9hdmlmICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgYXZpZjsKICAgIGltYWdlL3BuZyAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBwbmc7CiAgICBpbWFnZS9zdmcreG1sICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgc3ZnIHN2Z3o7CiAgICBpbWFnZS90aWZmICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgdGlmIHRpZmY7CiAgICBpbWFnZS92bmQud2FwLndibXAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgd2JtcDsKICAgIGltYWdlL3dlYnAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB3ZWJwOwogICAgaW1hZ2UveC1pY29uICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIGljbzsKICAgIGltYWdlL3gtam5nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBqbmc7CiAgICBpbWFnZS94LW1zLWJtcCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgYm1wOwoKICAgIGZvbnQvd29mZiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB3b2ZmOwogICAgZm9udC93b2ZmMiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHdvZmYyOwoKICAgIGFwcGxpY2F0aW9uL2phdmEtYXJjaGl2ZSAgICAgICAgICAgICAgICAgICAgICAgICBqYXIgd2FyIGVhcjsKICAgIGFwcGxpY2F0aW9uL2pzb24gICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBqc29uOwogICAgYXBwbGljYXRpb24vbWFjLWJpbmhleDQwICAgICAgICAgICAgICAgICAgICAgICAgIGhxeDsKICAgIGFwcGxpY2F0aW9uL21zd29yZCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBkb2M7CiAgICBhcHBsaWNhdGlvbi9wZGYgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgcGRmOwogICAgYXBwbGljYXRpb24vcG9zdHNjcmlwdCAgICAgICAgICAgICAgICAgICAgICAgICAgIHBzIGVwcyBhaTsKICAgIGFwcGxpY2F0aW9uL3J0ZiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBydGY7CiAgICBhcHBsaWNhdGlvbi92bmQuYXBwbGUubXBlZ3VybCAgICAgICAgICAgICAgICAgICAgbTN1ODsKICAgIGFwcGxpY2F0aW9uL3ZuZC5nb29nbGUtZWFydGgua21sK3htbCAgICAgICAgICAgICBrbWw7CiAgICBhcHBsaWNhdGlvbi92bmQuZ29vZ2xlLWVhcnRoLmtteiAgICAgICAgICAgICAgICAga216OwogICAgYXBwbGljYXRpb24vdm5kLm1zLWV4Y2VsICAgICAgICAgICAgICAgICAgICAgICAgIHhsczsKICAgIGFwcGxpY2F0aW9uL3ZuZC5tcy1mb250b2JqZWN0ICAgICAgICAgICAgICAgICAgICBlb3Q7CiAgICBhcHBsaWNhdGlvbi92bmQubXMtcG93ZXJwb2ludCAgICAgICAgICAgICAgICAgICAgcHB0OwogICAgYXBwbGljYXRpb24vdm5kLm9hc2lzLm9wZW5kb2N1bWVudC5ncmFwaGljcyAgICAgIG9kZzsKICAgIGFwcGxpY2F0aW9uL3ZuZC5vYXNpcy5vcGVuZG9jdW1lbnQucHJlc2VudGF0aW9uICBvZHA7CiAgICBhcHBsaWNhdGlvbi92bmQub2FzaXMub3BlbmRvY3VtZW50LnNwcmVhZHNoZWV0ICAgb2RzOwogICAgYXBwbGljYXRpb24vdm5kLm9hc2lzLm9wZW5kb2N1bWVudC50ZXh0ICAgICAgICAgIG9kdDsKICAgIGFwcGxpY2F0aW9uL3ZuZC5vcGVueG1sZm9ybWF0cy1vZmZpY2Vkb2N1bWVudC5wcmVzZW50YXRpb25tbC5wcmVzZW50YXRpb24KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBwcHR4OwogICAgYXBwbGljYXRpb24vdm5kLm9wZW54bWxmb3JtYXRzLW9mZmljZWRvY3VtZW50LnNwcmVhZHNoZWV0bWwuc2hlZXQKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB4bHN4OwogICAgYXBwbGljYXRpb24vdm5kLm9wZW54bWxmb3JtYXRzLW9mZmljZWRvY3VtZW50LndvcmRwcm9jZXNzaW5nbWwuZG9jdW1lbnQKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBkb2N4OwogICAgYXBwbGljYXRpb24vdm5kLndhcC53bWxjICAgICAgICAgICAgICAgICAgICAgICAgIHdtbGM7CiAgICBhcHBsaWNhdGlvbi93YXNtICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgd2FzbTsKICAgIGFwcGxpY2F0aW9uL3gtN3otY29tcHJlc3NlZCAgICAgICAgICAgICAgICAgICAgICA3ejsKICAgIGFwcGxpY2F0aW9uL3gtY29jb2EgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBjY287CiAgICBhcHBsaWNhdGlvbi94LWphdmEtYXJjaGl2ZS1kaWZmICAgICAgICAgICAgICAgICAgamFyZGlmZjsKICAgIGFwcGxpY2F0aW9uL3gtamF2YS1qbmxwLWZpbGUgICAgICAgICAgICAgICAgICAgICBqbmxwOwogICAgYXBwbGljYXRpb24veC1tYWtlc2VsZiAgICAgICAgICAgICAgICAgICAgICAgICAgIHJ1bjsKICAgIGFwcGxpY2F0aW9uL3gtcGVybCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBwbCBwbTsKICAgIGFwcGxpY2F0aW9uL3gtcGlsb3QgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBwcmMgcGRiOwogICAgYXBwbGljYXRpb24veC1yYXItY29tcHJlc3NlZCAgICAgICAgICAgICAgICAgICAgIHJhcjsKICAgIGFwcGxpY2F0aW9uL3gtcmVkaGF0LXBhY2thZ2UtbWFuYWdlciAgICAgICAgICAgICBycG07CiAgICBhcHBsaWNhdGlvbi94LXNlYSAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgc2VhOwogICAgYXBwbGljYXRpb24veC1zaG9ja3dhdmUtZmxhc2ggICAgICAgICAgICAgICAgICAgIHN3ZjsKICAgIGFwcGxpY2F0aW9uL3gtc3R1ZmZpdCAgICAgICAgICAgICAgICAgICAgICAgICAgICBzaXQ7CiAgICBhcHBsaWNhdGlvbi94LXRjbCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgdGNsIHRrOwogICAgYXBwbGljYXRpb24veC14NTA5LWNhLWNlcnQgICAgICAgICAgICAgICAgICAgICAgIGRlciBwZW0gY3J0OwogICAgYXBwbGljYXRpb24veC14cGluc3RhbGwgICAgICAgICAgICAgICAgICAgICAgICAgIHhwaTsKICAgIGFwcGxpY2F0aW9uL3hodG1sK3htbCAgICAgICAgICAgICAgICAgICAgICAgICAgICB4aHRtbDsKICAgIGFwcGxpY2F0aW9uL3hzcGYreG1sICAgICAgICAgICAgICAgICAgICAgICAgICAgICB4c3BmOwogICAgYXBwbGljYXRpb24vemlwICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHppcDsKCiAgICBhcHBsaWNhdGlvbi9vY3RldC1zdHJlYW0gICAgICAgICAgICAgICAgICAgICAgICAgYmluIGV4ZSBkbGw7CiAgICBhcHBsaWNhdGlvbi9vY3RldC1zdHJlYW0gICAgICAgICAgICAgICAgICAgICAgICAgZGViOwogICAgYXBwbGljYXRpb24vb2N0ZXQtc3RyZWFtICAgICAgICAgICAgICAgICAgICAgICAgIGRtZzsKICAgIGFwcGxpY2F0aW9uL29jdGV0LXN0cmVhbSAgICAgICAgICAgICAgICAgICAgICAgICBpc28gaW1nOwogICAgYXBwbGljYXRpb24vb2N0ZXQtc3RyZWFtICAgICAgICAgICAgICAgICAgICAgICAgIG1zaSBtc3AgbXNtOwoKICAgIGF1ZGlvL21pZGkgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBtaWQgbWlkaSBrYXI7CiAgICBhdWRpby9tcGVnICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgbXAzOwogICAgYXVkaW8vb2dnICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIG9nZzsKICAgIGF1ZGlvL3gtbTRhICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBtNGE7CiAgICBhdWRpby94LXJlYWxhdWRpbyAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgcmE7CgogICAgdmlkZW8vM2dwcCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDNncHAgM2dwOwogICAgdmlkZW8vbXAydCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHRzOwogICAgdmlkZW8vbXA0ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIG1wNDsKICAgIHZpZGVvL21wZWcgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBtcGVnIG1wZzsKICAgIHZpZGVvL3F1aWNrdGltZSAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBtb3Y7CiAgICB2aWRlby93ZWJtICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgd2VibTsKICAgIHZpZGVvL3gtZmx2ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBmbHY7CiAgICB2aWRlby94LW00diAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgbTR2OwogICAgdmlkZW8veC1tbmcgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIG1uZzsKICAgIHZpZGVvL3gtbXMtYXNmICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBhc3ggYXNmOwogICAgdmlkZW8veC1tcy13bXYgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHdtdjsKICAgIHZpZGVvL3gtbXN2aWRlbyAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBhdmk7Cn0K",
                    Name: "/etc/nginx/mime.types",
                },
                {
                    Contents: "c2VydmVyew0KICAgbGlzdGVuICAgODA4MDsgDQogICBzZXJ2ZXJfbmFtZSBleGFtcGxlLmNvbTsNCiAgIHN0YXR1c196b25lIGV4YW1wbGUuY29tX2Rtel9uZ2lueDsNCiANCiAgIGxvY2F0aW9uIC9mb28gew0KICAgICAgICBzdGF0dXNfem9uZSBleGFtcGxlLmNvbV9kbXpfbmdpbngtZm9vOw0KICAgICAgICBkZWZhdWx0X3R5cGUgInRleHQvcGxhaW4iOw0KICAgICAgICByZXR1cm4gMjAwICdPS1xuJzsNCiAgICB9DQoNCiAgICBsb2NhdGlvbiAvYmFyIHsNCiAgICAgICAgc3RhdHVzX3pvbmUgZXhhbXBsZS5jb21fZG16X25naW54LWJhcjsNCiAgICAgICAgZGVmYXVsdF90eXBlICJ0ZXh0L3BsYWluIjsNCiAgICAgICAgcmV0dXJuIDIwMCAnT0tcbic7DQogICAgfQ0KDQp9",
                    Name: "/etc/nginx/conf.d/app.conf",
                },
                {
                    Contents: "c2VydmVyIHsNCiAgICBsaXN0ZW4gODAwMTsNCiAgICByb290IC91c3Ivc2hhcmUvbmdpbngvaHRtbDsNCiAgICBhY2Nlc3NfbG9nIG9mZjsNCiAgICBsb2NhdGlvbiAgPSAvZGFzaGJvYXJkLmh0bWwgew0KICAgIH0NCiAgICBhbGxvdyAwLjAuMC4wLzA7DQogICAgZGVueSBhbGw7DQogICAgbG9jYXRpb24gL2FwaSB7DQogICAgICAgIGFwaSB3cml0ZT1vbjsNCiAgICB9DQp9DQojIyMK",
                    Name: "/etc/nginx/conf.d/api.conf",
                },
            },
            RootDir: "/etc/nginx",
        }

        directoryMap := DirectoryMap{
		"/etc/nginx": Directory{
			Files: []FileData{
				{
					Contents: "",
					Name:     "/etc/nginx/conf.d/api.conf",
					Size:     219,
				},
				{
					Contents: "",
					Name:     "/etc/nginx/conf.d/app.conf",
					Size:     384,
				},
				{
					Contents: "",
					Name:     "/etc/nginx/mime.types",
					Size:     5349,
				},
				{
					Contents: "",
					Name:     "/etc/nginx/nginx.conf",
					Size:     1029,
				},
			},
			Name:        "/etc/nginx",
			Permissions: "",
			Size:        0,
			UpdateTime:  time.Time{},
		},
	}

        instances := []string{
		"b369ed9a-e478-58ad-9c78-6363f93d5aaf",
		"52e29895-addf-545a-8c03-32a05cf58251",
	}
    
        response := InstanceGroupConfigResponse{
            AuxFiles:       auxFiles,
            ConfigFiles:    configFiles,
            CreateTime:     time.Now(),
            DirectoryMap:   directoryMap,
            Instances:      instances,
            UpdateTime:     time.Now(),
        }
        jsonResponse, err := json.Marshal(response)
        if err != nil {
            http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResponse)
    case http.MethodPost:
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

        selfLinks := SelfLinks{
            Rel: "/api/platform/v1/instance-groups/deployments/09aa57d6-bff8-437e-9796-4f5308c454b0",
        }
        publishConfigResponse := PublishConfigResponse{
            DeploymentUID: "09aa57d6-bff8-437e-9796-4f5308c454b0",
            Links:         selfLinks,
            Result:        "Publish group configuration request Accepted",
        }

        jsonResponse, err := json.Marshal(publishConfigResponse)
        if err != nil {
            http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResponse)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
