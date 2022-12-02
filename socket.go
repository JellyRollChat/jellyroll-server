package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type CommPacket struct {
}

func SocketAPI(keyCollection *ED25519Keys) {

	api := mux.NewRouter()

	// Channel Socket
	api.HandleFunc("/talk", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] [%s] +client\n"+white, timeStamp(), conn.RemoteAddr())
		socketParser(conn, keyCollection)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(clientCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func socketParser(conn *websocket.Conn, keyCollection *ED25519Keys) {

	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// socket session closed
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
			break
		}

		msgType := 1

		if bytes.HasPrefix(msg, pingMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		} else if bytes.HasPrefix(msg, mesgMsg) {
			fmt.Println("Full name:")
			fullnameb := hex.EncodeToString([]byte("username@server.tld"))
			fmt.Println(fullnameb)
			fmt.Println("Short name: ")
			shortnameb := hex.EncodeToString([]byte("username"))
			fmt.Println(shortnameb)
			fmt.Println("Body: ")
			bodyb := hex.EncodeToString([]byte("This is my first mesg"))
			fmt.Println(bodyb)
			fmt.Println("Full Message: ")
			fullmesgb := "<<to:" + fullnameb + "::body:" + bodyb + ">>"
			fmt.Println(fullmesgb)
			// message format
			// <<to:username@server.tld::body:54686973206973206d79206669727374206d657367>>

			conn.WriteMessage(msgType, []byte(fullmesgb))

			thisDecodedMesg, _ := hex.DecodeString(fullmesgb)

			conn.WriteMessage(msgType, thisDecodedMesg)
		}

	}
}
