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
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
			break
		}

		thisMessage := Message{}

		unmarshalError := json.Unmarshal(msg, &thisMessage)
		if unmarshalError != nil {
			log.Println(unmarshalError)
		}

		triageSocketMsg(&thisMessage, conn)
		conn.WriteJSON(&thisMessage)

	}

}

func triageSocketMsg(msg *Message, conn *websocket.Conn) {

	if msg.Type == 100 {
		log.Println("message type 100")
		log.Println("this is a friend request")
		conn.WriteMessage(1, []byte("friend request"))
	} else if msg.Type == 200 {
		log.Println("message type 200")
		log.Println("this is a normal chat message")
		conn.WriteMessage(1, []byte("chat message"))
	} else {
		log.Println("I didn't understand this message")
		conn.WriteMessage(1, []byte("didn't understand this request"))
	}

}
