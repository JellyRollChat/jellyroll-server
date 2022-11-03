package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func TrackerSocketAPI(keyCollection *ED25519Keys) {

	api := mux.NewRouter()

	// Channel Socket
	api.HandleFunc("/channel/tracker", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] [%s] +tracker\n"+white, timeStamp(), conn.RemoteAddr())
		socketTrackerParser(conn, keyCollection)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(trackerCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func socketTrackerParser(conn *websocket.Conn, keyCollection *ED25519Keys) {

	// Start reading messages in a loop
	for {

		// When we're done close the connection.
		defer conn.Close()

		// Try to read the message, if there's an error, shit the bed.
		msgType, msg, err := conn.ReadMessage()
		if err != nil {

			// socket session closed
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())

			return
		}

		// If we recognize this message as a tracker ping, run this.
		if bytes.HasPrefix(msg, pingMsg) {

			// reply to PING with PONG
			conn.WriteMessage(msgType, pongMsg)

			// print status message
			fmt.Printf(nc+"\n[%s] [%s] PING!\n"+white, timeStamp(), conn.RemoteAddr())

		}

	}
}
