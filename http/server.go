package main

import ( 
    "fmt"
    "net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
    fmt.Fprintf(w, "SUCCESS\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
    
    fmt.Fprintf(w, "client: %v\n",  req.RemoteAddr)

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func main() {
    http.HandleFunc("/", hello)
    http.HandleFunc("/headers", headers)
    http.ListenAndServe("0.0.0.0:8080", nil)
}
