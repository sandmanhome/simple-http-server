package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func GetQueryHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	age := mux.Vars(r)["age"]
	json.NewEncoder(w).Encode(map[string]string{"name": name, "age": age})
}
