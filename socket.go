package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
		socketHandler(conn, keyCollection)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(clientCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func socketHandler(conn *websocket.Conn, keyCollection *ED25519Keys) {

	defer conn.Close()
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
		return
	}

	thisMessage := Packet{}

	unmarshalError := json.Unmarshal(msg, &thisMessage)
	if unmarshalError != nil {
		log.Println("unmarshal error", unmarshalError)
	}

	socketMsgRouter(&thisMessage, conn)
	conn.WriteJSON(&thisMessage)

}

func socketMsgRouter(msg *Packet, conn *websocket.Conn) {

	if msg.MsgType == 100 {
		log.Println("message type 100")
		log.Println("this is a login auth request")
		loginHandler(msg, conn)
	} else {
		log.Println("I didn't understand this message")
		conn.WriteMessage(1, []byte("didn't understand this request"))
	}

}

func authdSocketMsgWriter(conn *websocket.Conn) {

	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
			return
		}

		thisMessage := Packet{}

		unmarshalError := json.Unmarshal(msg, &thisMessage)
		if unmarshalError != nil {
			log.Println("unmarshal error", unmarshalError)
		}

		if thisMessage.MsgType == 100 {
			log.Println("100: User is already logged in.")
			conn.WriteMessage(1, []byte("You are already logged in."))
		} else if thisMessage.MsgType == 200 {
			log.Println("200: this is a friend request")
			conn.WriteMessage(1, []byte("friend request"))
		} else if thisMessage.MsgType == 300 {
			log.Println("300: this is a normal chat message")
			conn.WriteMessage(1, []byte("chat message"))
		} else {
			log.Println("???: I didn't understand this message")
			conn.WriteMessage(1, []byte("didn't understand this request"))
		}

	}

}

func splitUserPassStr(userpass string) []string {
	return strings.Split(userpass, ",")
}

func loginHandler(msg *Packet, conn *websocket.Conn) {
	log.Println("loginHandler reached")
	log.Println("Login contents: ", msg.MsgContent)
	userpass := splitUserPassStr(msg.MsgContent)
	log.Println("Login username: ", userpass[0])
	log.Println("Login password: ", userpass[1])
	if stringExistsInFile(msg.MsgContent) {
		log.Println("User exists in user list")
		conn.WriteMessage(1, []byte("Welcome :)"))
		authdSocketMsgWriter(conn)
	} else {
		log.Println("User does not exist in user list")
		conn.WriteMessage(1, []byte("Access Denied. Goodbye!"))
		conn.Close()
	}

}

func stringExistsInFile(thisString string) bool {
	f, err := os.Open("admin/users.list")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if thisString == scanner.Text() {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return false
}
