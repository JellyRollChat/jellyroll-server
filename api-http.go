package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func serverWebAPI() {

	// Init API
	api := mux.NewRouter()
	// api := r.PathPrefix("/api/v1").Subrouter()

	// Home
	api.HandleFunc("/", home).Methods(http.MethodGet)

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(webPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("", w, r)
	w.Write([]byte("welcome"))
}
func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf(brightgreen+"\n/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+brightcyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}
