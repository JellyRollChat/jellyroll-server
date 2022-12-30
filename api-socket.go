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
	announce("Socket API")
	api := mux.NewRouter()
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
		"*",
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

	// Channel Socket
	api.HandleFunc("/talk", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] +client [%s]\n"+nc, timeStamp(), conn.RemoteAddr())
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

	if thisMessage.Type != 100 {
		log.Println("User attempted message without valid login session, disconnecting.")
		conn.WriteMessage(1, []byte("Access Denied!"))
		conn.Close()
	}
	// socketMsgRouter(&thisMessage, conn)
	loginHandler(&thisMessage, conn)
	conn.WriteJSON(&thisMessage)

}

func loginHandler(msg *Packet, conn *websocket.Conn) {
	userpass := splitUserPassStr(msg.Content)
	log.Println("Login username: ", userpass[0])
	if stringExistsInFile(msg.Content) {
		thisSession := UserSession{
			Username:   userpass[0] + "@" + userpass[1],
			State:      ClientStateExchange{},
			Conn:       conn,
			Authorized: true,
		}
		if !fileExists("admin/users/" + userpass[0] + ".state") {
			createFile("admin/users/" + userpass[0] + ".state")
			thisSession.State = ClientStateExchange{
				CurrentFriends: []string{
					"esp@3ck0.com",
					"fred@server.3ck0.com",
				},
				PendingFriends: []string{},
				BlockedFriends: []string{},
				BlockedServers: []string{},
			}
			for _, r := range thisSession.State.CurrentFriends {
				log.Println("Online? ", r)
				for _, rr := range GlobalUserSessions {
					log.Println("Checking: ", rr)
					if rr.Username == r {
						r = "(" + r + ")"
						log.Println("Match found! ", r)
					}
				}
			}
			marshalState, msErr := json.Marshal(thisSession.State)
			if msErr != nil {
				log.Println("Marshal Error: " + msErr.Error())
				conn.WriteMessage(1, []byte("Marshal Error"))
				conn.Close()
				return
			}
			writeFile("admin/users/"+userpass[0]+".state", string(marshalState))
			// thisFile := readFile("admin/users/" + userpass[0] + ".state")
			// log.Println("This user's state: " + thisFile)
			conn.WriteJSON(marshalState)
			// log.Println("User state created: " + "admin/users/" + userpass[0] + ".state")

		} else if stringExistsInFile(msg.Content) {
			// log.Println("User state exists: " + "admin/users/" + userpass[0] + ".state")

			thisFile := readFile("admin/users/" + userpass[0] + ".state")
			unmarshErr := json.Unmarshal([]byte(thisFile), &thisSession.State)
			for _, r := range thisSession.State.CurrentFriends {
				log.Println("Online? ", r)
				for _, rr := range GlobalUserSessions {
					log.Println("Checking: ", rr)
					if rr.Username == r {
						r = "(" + r + ")"
						log.Println("Match found! ", r)
					}
				}
			}
			if unmarshErr != nil {
				log.Println("Unmarshal error: " + unmarshErr.Error())
			}

			conn.WriteMessage(1, []byte(thisFile))
		}
		// log.Println(msg.Content)
		// log.Println("User exists in user list")
		AddUserSession(&thisSession)
		log.Println(brightcyan+"Global Socket Sessions: ", len(GlobalUserSessions))
		log.Println("Users online: ")
		for i, session := range GlobalUserSessions {
			fmt.Printf(i + " " + session.Username + " ")
		}

		authdSocketMsgWriter(conn)
	} else {
		log.Println("User does not exist in user list")
		conn.WriteMessage(1, []byte("Access Denied. Goodbye!"))
		conn.Close()
	}

}

func (s *UserSession) Listen() {
	mutex.Lock()
	defer mutex.Unlock()
	defer RemoveUserSession(s.Username)
	if s.Conn != nil {
		return
	}
	for {
		var msg ClientMessage
		err := s.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("error reading message:", err)
			break
		}
		var packet Packet
		err = json.Unmarshal([]byte(msg.Body), &packet)
		if err != nil {
			log.Println("error unmarshalling message body:", err)
			continue
		}
		log.Println(packet)
		switch packet.Type {
		case 100:
			if s.Authorized {
				log.Println("100: this is a login request.")
				log.Println("User is already logged in")
			}
		case 200:
			log.Println("200: this is a friend request")
		case 300:
			log.Println("300: this is a normal chat message")
		default:
			log.Println("???: I didn't understand this message")
		}
	}
}

func authdSocketMsgWriter(conn *websocket.Conn) {

	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
			log.Println(brightcyan+"Global Socket Sessions: ", len(GlobalUserSessions))

			return
		}

		thisPacket := Packet{}

		unmarshalError := json.Unmarshal(msg, &thisPacket)
		if unmarshalError != nil {
			log.Println("unmarshal error", unmarshalError)
		}

		if thisPacket.Type == 100 {
			log.Println("100: User is already logged in.")
			conn.WriteMessage(1, []byte("You are already logged in."))
		} else if thisPacket.Type == 200 {
			log.Println("200: this is a friend request")
			conn.WriteMessage(1, []byte("friend request"))
		} else if thisPacket.Type == 300 {
			log.Println("300: this is a normal chat message")
			conn.WriteMessage(1, []byte("chat message"))
		} else {
			log.Println("???: I didn't understand this message")
			conn.WriteMessage(1, []byte("didn't understand this request"))
		}

	}

}

// AddUserSession adds a new user session to the global map
func AddUserSession(s *UserSession) {
	GlobalUserSessions[s.Username] = s
	go s.Listen()
}

// RemoveUserSession removes a user session from the global map
func RemoveUserSession(username string) {
	delete(GlobalUserSessions, username)
}
