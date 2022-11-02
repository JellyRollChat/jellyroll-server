package main

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func TrackerSocketAPI(keyCollection *ED25519Keys) {

	api := mux.NewRouter()

	// Channel Socket
	api.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
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

	for {

		defer conn.Close()

		msgType, msg, err := conn.ReadMessage()

		if err != nil {

			// socket session closed
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())

			return
		}

		// if the socket has a join message
		if bytes.HasPrefix(msg, joinMsg) {

			msgToString := string(msg)

			trimNewlinePrefix := strings.TrimRight(msgToString, "\n")

			trimmedPubKey := strings.TrimLeft(trimNewlinePrefix, "JOIN ")

			// complains about loop duration
			regValidate, _ := regexp.MatchString(`[a-f0-9]{64}`, trimmedPubKey[:64])
			if !regValidate {
				fmt.Printf("\nContains illegal characters")
				conn.Close()
				return
			}

		}
		if bytes.HasPrefix(msg, rtrnMsg) {
			fmt.Printf("\nreturn message: %s", string(msg))

			input := strings.TrimLeft(string(msg), "RTRN ")
			var cert = strings.Split(input, " ")

			if !verifySignature(cert[0], keyCollection.publicKey, cert[1]) {
				fmt.Printf("\nsig doesnt verify")
			}
			if verifySignature(cert[0], keyCollection.publicKey, cert[1]) {
				fmt.Printf("\nsig verifies")
			}
		}
		if bytes.HasPrefix(msg, pubkMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}
		if bytes.HasPrefix(msg, peerMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}

	}
}
