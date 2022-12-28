package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func WebAPI() {

	corsAllowedHeaders := []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Cache-Control",
		"Content-Security-Policy",
		"Feature-Policy",
		"Referrer-Policy",
		"X-Requested-With"}

	corsOrigins := []string{
		"http://127.0.0.1:1430",
	}

	corsMethods := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"OPTIONS",
	}

	headersCORS := handlers.AllowedHeaders(corsAllowedHeaders)
	originsCORS := handlers.AllowedOrigins(corsOrigins)
	methodsCORS := handlers.AllowedMethods(corsMethods)

	log.Println("API launched")
	api := mux.NewRouter()

	api.HandleFunc("/signup", SignupHandlerPOST).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/signup/", SignupHandlerPOST).Methods(http.MethodPost, http.MethodOptions)

	http.ListenAndServe(":"+strconv.Itoa(webPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf("\n"+purple+r.Method+" "+brightgreen+"/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+cyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

func SignupHandlerPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	log.Println("Request headers:", r.Header)
	log.Println("Request form:", r.Form)
	reportRequest("signup", w, r)
	thisSignup := AuthObject{}
	r.ParseForm()
	for key, value := range r.Form {
		log.Println("key: ", key, "value: ", value)
		if key == "username" {
			thisSignup.Username = sanitizeString(value[0], 20)
		} else if key == "signupPassword" {
			thisSignup.Password = value[0]
		}
	}
	err := json.NewDecoder(r.Body).Decode(&thisSignup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !fileContainsString(thisSignup.Username, "admin/users.list") {
		thisSignup.Password = hashit(thisSignup.Password)
		appendFile("admin/users.list", thisSignup.Username+"@"+servertld+","+thisSignup.Password+"\n")
		type userInfo struct {
			Username  string
			Servertld string
		}
		thisUser := userInfo{}
		thisUser.Servertld = servertld
		thisUser.Username = thisSignup.Username
		fmt.Fprintf(w, "OK!")
		log.Println(brightmagenta + "New User: " + magenta + thisSignup.Username + "@" + servertld)
	} else {
		log.Println(brightred + "Username is not available " + nc)
		fmt.Fprintf(w, "DENIED!")
		return
	}
	log.Println("Response headers:", w.Header())
}
