package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// func httpTrackerAPI(keyCollection *ED25519Keys) {

// 	api := mux.NewRouter()

// 	// Channel Socket
// 	api.HandleFunc("/channel/tracker", func(w http.ResponseWriter, r *http.Request) {
// 		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 		conn, _ := upgrader.Upgrade(w, r, nil)
// 		defer conn.Close()
// 		fmt.Printf(brightgreen+"\n[%s] [%s] +tracker\n"+white, timeStamp(), conn.RemoteAddr())
// 		httpTrackerHandler(conn, keyCollection)
// 	})

// 	// Serve via HTTP
// 	http.ListenAndServe(":"+strconv.Itoa(trackerCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
// }

func httpTrackerAPI() {

	api := mux.NewRouter()
	api.HandleFunc("/tracker/ping", trackerPingHandler).Methods(http.MethodGet)

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(trackerCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func trackerPingHandler(w http.ResponseWriter, r *http.Request) {

	// Set content type header for text
	w.Header().Set("Content-Type", "application/json")

	// Assemble this into the header
	w.WriteHeader(http.StatusOK)

	// Announce that someone has hit this endpoint
	reportRequest("PING", w, r)

	// Write the full response with header and serve to the user
	w.Write([]byte("PONG"))
}
