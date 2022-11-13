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
	api.HandleFunc("/signup", SignupHandlerPOST).Methods(http.MethodPost)
	api.HandleFunc("/signup/", SignupHandlerPOST).Methods(http.MethodPost)

	http.ListenAndServe(":"+strconv.Itoa(webPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf("\n"+purple+r.Method+" "+brightgreen+"/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+cyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

func SignupHandlerGET(w http.ResponseWriter, r *http.Request) {
	reportRequest("signup", w, r)
	files := []string{
		"templates/signup.html",
	}

	t, parseSignupFiles := template.ParseFiles(files...)

	handle("", parseSignupFiles)
	if parseSignupFiles != nil {

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	whatswrong := t.Execute(w, servertld)
	handle("http signup render error", whatswrong)

}

func SignupHandlerPOST(w http.ResponseWriter, r *http.Request) {
	reportRequest("signup", w, r)

	type signupObject struct {
		username   string
		password   string
		passRepeat string
	}

	thisSignup := signupObject{}

	r.ParseForm()
	for key, value := range r.Form {
		if key == "signupUsername" {
			thisSignup.username = sanitizeString(value[0], 20)
		} else if key == "signupPassword" {
			thisSignup.password = value[0]
		} else if key == "signupPasswordRepeat" {
			thisSignup.passRepeat = value[0]
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

		thisSignup.password = hashit(thisSignup.password)

		appendFile("admin/users.list", thisSignup.username+"::"+thisSignup.password+"\n")

		files := []string{
			"templates/signupSuccess.html",
		}

		t, parseSignupFiles := template.ParseFiles(files...)

		handle("error parsing signupsuccess view: ", parseSignupFiles)
		if parseSignupFiles != nil {

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		type userInfo struct {
			Username  string
			Servertld string
		}

		thisUser := userInfo{}

		thisUser.Servertld = servertld
		thisUser.Username = thisSignup.username

		whatswrong := t.Execute(w, thisUser)
		handle("http signup success render error", whatswrong)

		// fmt.Fprintf(w, "Success! You can now use your username and password to login. \n\nUsername: %s\nServer: %s\n\nGive your friends this address: %s\n\nTip: It looks like an email but it's really your full username.", thisSignup.username, servertld, thisSignup.username+"@"+servertld)

		log.Println(brightmagenta + "New User: " + magenta + thisSignup.username + "@" + servertld)
		// return

	} else {
		log.Println(brightred + "Username is not available " + nc)
		fmt.Fprintf(w, "Username is not available ")
		return
	}

}
