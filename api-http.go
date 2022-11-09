package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func serverWebAPI() {

	api := mux.NewRouter()
	api.HandleFunc("/signup", SignupHandlerGET).Methods(http.MethodGet)
	api.HandleFunc("/signup/", SignupHandlerGET).Methods(http.MethodGet)
	api.HandleFunc("/signup", SignupHandlerGET).Methods(http.MethodPost)
	api.HandleFunc("/signup/", SignupHandlerGET).Methods(http.MethodPost)

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

func SignupHandlerPOST(w http.ResponseWriter, r *http.Request) {

	type signupObject struct {
		username   string
		password   string
		passRepeat string
	}

	thisSignup := signupObject{}

	r.ParseForm()
	for key, value := range r.Form {
		if key == "signupUsername" {
			thisSignup.username = value[0]
		} else if key == "signupPassword" {
			thisSignup.password = value[0]
		} else if key == "signupPasswordRepeat" {
			thisSignup.passRepeat = value[0]
		} else {
			log.Println("login form key: " + key + "\nlogin form value: " + strings.Join(value, ", "))
		}
	}

	if !fileContainsString(thisSignup.username, "admin/users.list") {

		if len(thisSignup.username) < 6 {
			log.Println(brightred + "Username is too short " + nc)
			fmt.Fprintf(w, "Username is too short ")
			return
		}

		if len(thisSignup.username) > 20 {
			log.Println(brightred + "Username is too long " + nc)
			fmt.Fprintf(w, "Username is too long ")
			return
		}

		if strings.Compare(thisSignup.password, thisSignup.passRepeat) != 0 {
			log.Println(brightred + "Passwords don't match " + nc)
			fmt.Fprintf(w, "Passwords don't match ")
			return
		}

		if len(thisSignup.password) < 8 {
			log.Println(brightred + "password is too short " + nc)
			fmt.Fprintf(w, "Password is too short ")
			return
		}

		if len(thisSignup.password) > 120 {
			log.Println(brightred + "password is too long " + nc)
			fmt.Fprintf(w, "Password is too long ")
			return
		}

		writeFile("admin/users.list", thisSignup.username+"::"+hashit(thisSignup.password))

	} else {
		log.Println(brightred + "Username is not available " + nc)
		fmt.Fprintf(w, "Username is not available ")
		return
	}

}
