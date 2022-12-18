package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func serverWebAPI() {

	api := mux.NewRouter()
	api.HandleFunc("/signup", SignupHandlerPOST).Methods(http.MethodPost)
	api.HandleFunc("/signup/", SignupHandlerPOST).Methods(http.MethodPost)

	http.ListenAndServe(":"+strconv.Itoa(webPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf("\n"+purple+r.Method+" "+brightgreen+"/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+cyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

func SignupHandlerPOST(w http.ResponseWriter, r *http.Request) {
	reportRequest("signup", w, r)

	type SignupObject struct {
		Username string
		Password string
	}

	thisSignup := SignupObject{}

	// unmarshall json string to struct
	err := json.NewDecoder(r.Body).Decode(&thisSignup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.ParseForm()
	for key, value := range r.Form {
		if key == "signupUsername" {
			thisSignup.Username = sanitizeString(value[0], 20)
		} else if key == "signupPassword" {
			thisSignup.Password = value[0]
		}
	}

	if !fileContainsString(thisSignup.Username, "admin/users.list") {

		if len(thisSignup.Username) < 6 {
			log.Println(brightred + "Username is too short " + nc)
			fmt.Fprintf(w, "Username is too short ")
			return
		}

		if len(thisSignup.Username) > 20 {
			log.Println(brightred + "Username is too long " + nc)
			fmt.Fprintf(w, "Username is too long ")
			return
		}

		if len(thisSignup.Password) < 8 {
			log.Println(brightred + "password is too short " + nc)
			fmt.Fprintf(w, "Password is too short ")
			return
		}

		if len(thisSignup.Password) > 120 {
			log.Println(brightred + "password is too long " + nc)
			fmt.Fprintf(w, "Password is too long ")
			return
		}

		thisSignup.Password = hashit(thisSignup.Password)

		appendFile("admin/users.list", thisSignup.Username+"::"+thisSignup.Password+"\n")

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
		thisUser.Username = thisSignup.Username

		whatswrong := t.Execute(w, thisUser)
		handle("http signup success render error", whatswrong)

		fmt.Fprintf(w, "\"OK\"")

		log.Println(brightmagenta + "New User: " + magenta + thisSignup.Username + "@" + servertld)
		// return

	} else {
		log.Println(brightred + "Username is not available " + nc)
		fmt.Fprintf(w, "Username is not available ")
		return
	}

}
