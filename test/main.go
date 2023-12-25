package main

import (
   "fmt"
   "runtime"
   "github.com/golang/glog"
)

var version = "1.0"

func main() {
    fmt.Println("test main start")

    commitHash, commitTime, dirtyBuild := getBuildInfo()
    fmt.Printf("NGINX Ingress Controller Version=%v Commit=%v Date=%v DirtyState=%v Arch=%v/%v Go=%v\n", version, commitHash, commitTime, dirtyBuild, runtime.GOOS, runtime.GOARCH, runtime.Version())

    parseFlags()


    templateExecutor, templateExecutorV2 := createTemplateExecutors()
}

func createTemplateExecutors() (*version1.TemplateExecutor, *version2.TemplateExecutor) {

        nginxConfTemplatePath := "nginx.tmpl"
        nginxIngressTemplatePath := "nginx.ingress.tmpl"
        nginxVirtualServerTemplatePath := "nginx.virtualserver.tmpl"
        nginxTransportServerTemplatePath := "nginx.transportserver.tmpl"

        templateExecutor, err := version1.NewTemplateExecutor(nginxConfTemplatePath, nginxIngressTemplatePath)
        if err != nil {
                glog.Fatalf("Error creating TemplateExecutor: %v", err)
        }       
                
        templateExecutorV2, err := version2.NewTemplateExecutor(nginxVirtualServerTemplatePath, nginxTransportServerTemplatePath)
        if err != nil {
                glog.Fatalf("Error creating TemplateExecutorV2: %v", err)
        }
        
        return templateExecutor, templateExecutorV2

}
