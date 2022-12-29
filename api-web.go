package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func WebAPI() {

	corsAllowedHeaders := []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Cache-Control",
		"Content-Security-Policy",
		"Content-Type",
		"Feature-Policy",
		"Referrer-Policy",
		"X-Requested-With"}

	corsOrigins := []string{
		"http://127.0.0.1:1430",
		"http://server.3ck0.com:5270",
	}

	corsMethods := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"OPTIONS",
	}

	log.Println("API launched")

	http.HandleFunc("/signup", SignupHandlerPOST)
	http.HandleFunc("/signup/", SignupHandlerPOST)

	http.ListenAndServe(":"+strconv.Itoa(webPort), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, origin := range corsOrigins {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		for _, method := range corsMethods {
			w.Header().Add("Access-Control-Allow-Methods", method)
		}
		for _, header := range corsAllowedHeaders {
			w.Header().Add("Access-Control-Allow-Headers", header)
		}

		if r.Method == http.MethodOptions {
			enableCors(&w)
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:1430")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf("\n"+purple+r.Method+" "+brightgreen+"/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+cyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

func SignupHandlerPOST(w http.ResponseWriter, r *http.Request) {
	// Check the request method
	if r.Method == http.MethodOptions {
		// If the method is OPTIONS, return an ok
		w.WriteHeader(http.StatusOK)
		return
	}
	reportRequest("signup", w, r)
	parseerr := r.ParseForm()
	if parseerr != nil {
		log.Println("Form parse error on signup handler: ", parseerr)
	}
	thisSignup := AuthObject{}
	// for key, value := range r.Form {
	// 	log.Println("Key: ", key)
	// 	log.Println("value: ", value)
	// 	log.Println("key: ", key, "value: ", value)
	// 	if key == "username" {
	// 		thisSignup.Username = sanitizeString(value[0], 20)
	// 	} else if key == "signupPassword" {
	// 		thisSignup.Password = value[0]
	// 	}
	// }
	log.Println("Heres the body", r.Body)
	err := json.NewDecoder(r.Body).Decode(&thisSignup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !fileContainsString(thisSignup.Username, "admin/users.list") {
		// thisSignup.Password = hashit(thisSignup.Password)
		appendFile("admin/users.list", thisSignup.Username+"@"+servertld+","+thisSignup.Password+"\n")
		type userInfo struct {
			Username  string
			Servertld string
		}
		thisUser := userInfo{}
		thisUser.Servertld = servertld
		thisUser.Username = thisSignup.Username
		fmt.Fprintf(w, "\"OK\"")
		log.Println("New User: " + thisSignup.Username + "@" + servertld)
	} else {
		log.Println("Username is not available")
		fmt.Fprintf(w, "DENIED!")
		return
	}
	// log.Println("Response headers:", w.Header())
}
