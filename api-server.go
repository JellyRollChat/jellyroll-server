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

func ServerSocketAPI(keyCollection *ED25519Keys) {

	api := mux.NewRouter()

	// Channel Socket
	api.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] [%s] +server\n"+white, timeStamp(), conn.RemoteAddr())
		socketServerParser(conn, keyCollection)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(serverCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func socketServerParser(conn *websocket.Conn, keyCollection *ED25519Keys) {

	for {

		defer conn.Close()

		msgType, msg, err := conn.ReadMessage()

		if err != nil {

			// socket session closed
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())

			return
		}

		if bytes.HasPrefix(msg, peerMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}

	}
}
