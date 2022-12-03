package main

import (
	"encoding/json"
	"fmt"
	"log"
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

		thisMessage := Message{}
		unmarshalError := json.Unmarshal(msg, thisMessage)
		if unmarshalError != nil {
			log.Println(unmarshalError)
		}

		// msgType := 1

		conn.WriteJSON(thisMessage)
		// thisMessage := Message{
		// 	Type: 200,
		// 	From: "alex@server.3ck0.com",
		// 	Recv: "bess@server.3cko.com",
		// 	Body: "Hello this is a test",
		// }

		// thisMsgJson, thisMsgJsonErr := json.Marshal(thisMessage)
		// if thisMsgJsonErr != nil {
		// 	log.Println("There was an error marshalling the JSON for this message")
		// }

		// conn.WriteMessage(msgType, thisMsgJson)

	}

}
