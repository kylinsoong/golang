package mock

import (
    "net/http"
    "strconv"
    "encoding/json"
    "github.com/gorilla/mux"

    log  "github.com/kylinsoong/golang/nim-mock/pkg/vlogger"
)

type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

    log.Infof("Received POST request from %s for %s", r.RemoteAddr, r.URL.Path)

    vars := mux.Vars(r)
    userID := vars["id"]

    id, err := strconv.Atoi(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user := User{
        ID:   id,
	Name: "John Doe",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
