package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func serverWebAPI() {

	api := mux.NewRouter()
	api.HandleFunc("/signup", SignupHandlerGET).Methods(http.MethodGet)
	api.HandleFunc("/signup/", SignupHandlerGET).Methods(http.MethodGet)
	api.HandleFunc("/signup", SignupHandlerGET).Methods(http.MethodGet)
	api.HandleFunc("/signup/", SignupHandlerGET).Methods(http.MethodGet)

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(webPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf(brightgreen+"\n/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+brightcyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

func SignupHandlerGET(w http.ResponseWriter, r *http.Request) {

	files := []string{
		"templates/signup.html",
	}

	// Parse the file list
	t, parseSignupFiles := template.ParseFiles(files...)

	// if something goes wrong, report it, and where
	// the error was generated.
	handle("", parseSignupFiles)
	if parseSignupFiles != nil {

		// if something went wrong, browsers should be relayed a
		// an Internal Server Error status.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		// stop execution
		return
	}

	// if we did everything right, this should serve the request.
	whatswrong := t.Execute(w, r)
	handle("http signup render error", whatswrong)

}
