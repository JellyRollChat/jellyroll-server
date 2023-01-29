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
	if stringExistsInFile(msg.Content) {
		handleSuccessfulLogin(userpass, conn, msg)
	} else {
		conn.Close()
	}
}

// func loginHandler(msg *Packet, conn *websocket.Conn) {
// 	userpass := splitUserPassStr(msg.Content)
// 	log.Println("Login username: ", userpass[0])
// 	if stringExistsInFile(msg.Content) {
// 		thisSession := UserSession{
// 			Username:   userpass[0] + "@" + userpass[1],
// 			State:      ClientStateExchange{},
// 			Conn:       conn,
// 			Authorized: true,
// 		}
// 		if !fileExists("admin/users/" + userpass[0] + ".state") {
// 			createFile("admin/users/" + userpass[0] + ".state")
// 			thisSession.State = ClientStateExchange{
// 				CurrentFriends: []string{
// 					"esp@server.3ck0.com",
// 					"fred@server.3ck0.com",
// 				},
// 				PendingFriends: []string{
// 					"wanda@cool.yachts",
// 					"mark@white.monster",
// 				},
// 				BlockedFriends: []string{},
// 				BlockedServers: []string{},
// 			}
// 			for _, r := range thisSession.State.CurrentFriends {
// 				log.Println("Online? ", r)
// 				for _, rr := range GlobalUserSessions {
// 					log.Println("Checking: ", rr)
// 					if rr.Username == r {
// 						r = "(" + r + ")"
// 						log.Println("Match found! ", r)
// 					}
// 				}
// 			}
// 			marshalState, msErr := json.Marshal(thisSession.State)
// 			if msErr != nil {
// 				log.Println("Marshal Error: " + msErr.Error())
// 				conn.WriteMessage(1, []byte("Marshal Error"))
// 				conn.Close()
// 				return
// 			}
// 			writeFile("admin/users/"+userpass[0]+".state", string(marshalState))
// 			conn.WriteJSON(marshalState)
// 		} else if stringExistsInFile(msg.Content) {
// 			thisFile := readFile("admin/users/" + userpass[0] + ".state")
// 			unmarshErr := json.Unmarshal([]byte(thisFile), &thisSession.State)
// 			for _, r := range thisSession.State.CurrentFriends {
// 				log.Println("Online? ", r)
// 				for _, rr := range GlobalUserSessions {
// 					log.Println("Checking: ", rr)
// 					if rr.Username == r {
// 						r = "(" + r + ")"
// 						log.Println("Match found! ", r)
// 					}
// 				}
// 			}
// 			if unmarshErr != nil {
// 				log.Println("Unmarshal error: " + unmarshErr.Error())
// 			}
// 			conn.WriteMessage(1, []byte(thisFile))
// 		}
// 		AddUserSession(&thisSession)
// 		log.Println(brightcyan+"Global Socket Sessions: ", len(GlobalUserSessions))
// 		log.Println("Users online: ")
// 		for i, session := range GlobalUserSessions {
// 			fmt.Printf(i + " " + session.Username + " ")
// 		}
// 		authdSocketMsgWriter(conn)
// 	} else {
// 		log.Println("User does not exist in user list")
// 		conn.WriteMessage(1, []byte("Access Denied. Goodbye!"))
// 		conn.Close()
// 	}
// }

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
		switch packet.Type {
		case 100:
			if s.Authorized {
				log.Println("100: this is a login request.")
				log.Println("User is already logged in")
			}
		case 200:
			handleFriendRequest(msg, s)
		case 201:
			handleFriendApproval(msg, s)
		case 202:
			handleFriendDenial(msg, s)
		case 300:
			handleChatMessage(msg, s)
		default:
			log.Println("???: unhandled output")
		}

	}

}

func handleChatMessage(msg ClientMessage, s *UserSession) {

	// Split the sender and recipient addresses into their username and server parts
	_, servertld := splitAddress(msg.From)
	_, recvServerURL := splitAddress(msg.Recv)

	// Log the values of msg.From and msg.Recv
	log.Println("msg.From:", msg.From)
	log.Println("msg.Recv:", msg.Recv)

	// Log the values of servertld and recvServerURL
	log.Println("servertld:", servertld)
	log.Println("recvServerURL:", recvServerURL)

	// Log the contents of the GlobalFedServers map
	log.Println("GlobalFedServers:", GlobalFedServers)

	// Split the sender and recipient addresses into their username and server parts
	_, splitFrom := splitAddress(msg.From)
	_, splitTo := splitAddress(msg.Recv)

	// Log the values returned by the splitAddress function
	log.Println("splitAddress(msg.From):", splitFrom)
	log.Println("splitAddress(msg.Recv):", splitTo)

	// Check if the recipient is on a different server
	if recvServerURL != servertld {
		// Look up the federated server object
		f, ok := GlobalFedServers[recvServerURL]
		if !ok {
			// If the federated server is not registered, try to establish a websocket connection with it
			ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", recvServerURL), nil)
			if err != nil {
				log.Println("error connecting to federated server:", err)
				return
			}
			// Create a new FedServer object
			f := &FedServer{
				URL:       recvServerURL,
				Inbox:     make(chan *FedMessage),
				Outbox:    make(chan *FedMessage),
				Messages:  make(map[string]*FedMessage),
				Websocket: ws,
			}
			// Add the FedServer object to the GlobalFedServers map
			GlobalFedServers[recvServerURL] = f
		}

		// Start the FedServer's listen and send goroutines
		go f.Listen()
		go f.Send()
	}
}

func handleFriendDenial(msg ClientMessage, s *UserSession) {
	log.Println("202: this is a friend denial")
	// Check if the user has the pending friend in their list
	var pendingIndex int
	found := false
	for i, pendingFriend := range s.State.PendingFriends {
		if pendingFriend == msg.From {
			found = true
			pendingIndex = i
			break
		}
	}
	if !found {
		log.Println("error: user does not have pending friend in their list")
		return
	}

	// Remove the pending friend from the pending friends list
	s.State.PendingFriends = append(s.State.PendingFriends[:pendingIndex], s.State.PendingFriends[pendingIndex+1:]...)
	s.History = append(s.History, MessageHistoryEntry{
		Timestamp: msg.Timestamp,
		Message:   fmt.Sprintf("Friend request from %s denied", msg.From),
	})

	// Update the recipient's message history
	recipientSession, ok := GlobalUserSessions[msg.Recv]
	if ok {
		recipientSession.History = append(recipientSession.History, MessageHistoryEntry{
			Timestamp: msg.Timestamp,
			Message:   fmt.Sprintf("Friend request to %s denied", msg.Recv),
		})
	}
	// Store the updated user state in the database
	storeUserStateJSON(s.Username, s.State)
}

func handleFriendApproval(msg ClientMessage, s *UserSession) {
	// Check if the user has the pending friend in their list
	var pendingIndex int
	found := false
	for i, pendingFriend := range s.State.PendingFriends {
		if pendingFriend == msg.From {
			found = true
			pendingIndex = i
			break
		}
	}
	if !found {
		log.Println("error: user does not have pending friend in their list")
		return
	}
	// Remove the pending friend from the pending friends list
	s.State.PendingFriends = append(s.State.PendingFriends[:pendingIndex], s.State.PendingFriends[pendingIndex+1:]...)
	// Add the approved friend to the current friends list
	s.State.CurrentFriends = append(s.State.CurrentFriends, msg.From)
	s.History = append(s.History, MessageHistoryEntry{
		Timestamp: msg.Timestamp,
		Message:   fmt.Sprintf("Friend request from %s approved", msg.From),
	})

	// Update the recipient's message history
	recipientSession, ok := GlobalUserSessions[msg.Recv]
	if ok {
		recipientSession.State.CurrentFriends = append(recipientSession.State.CurrentFriends, msg.From)
		recipientSession.History = append(recipientSession.History, MessageHistoryEntry{
			Timestamp: msg.Timestamp,
			Message:   fmt.Sprintf("Friend request to %s approved", msg.Recv),
		})
	}
	// Store the updated user state in the database
	storeUserStateJSON(s.Username, s.State)
}

func handleFriendRequest(msg ClientMessage, s *UserSession) {
	log.Println("200: this is a friend request")
	s.State.PendingFriends = append(s.State.PendingFriends, msg.From)
	s.History = append(s.History, MessageHistoryEntry{
		Timestamp: msg.Timestamp,
		Message:   fmt.Sprintf("Friend request from %s", msg.From),
	})
	// Update the recipient's message history
	recipientSession, ok := GlobalUserSessions[msg.Recv]
	if ok {
		recipientSession.State.PendingFriends = append(recipientSession.State.PendingFriends, msg.From)
		recipientSession.History = append(recipientSession.History, MessageHistoryEntry{
			Timestamp: msg.Timestamp,
			Message:   fmt.Sprintf("Friend request from %s", msg.From),
		})
	}
	// Store the request in the database
	storeUserStateJSON(s.Username, s.State)
}

// func authdSocketMsgWriter(conn *websocket.Conn) {
// 	for {
// 		socketCount++
// 		userCount = countLines()
// 		defer conn.Close()
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
// 			socketCount--
// 			log.Println(brightcyan+"Global Socket Sessions: ", len(GlobalUserSessions))
// 			return
// 		}
// 		thisPacket := Packet{}
// 		unmarshalError := json.Unmarshal(msg, &thisPacket)
// 		if unmarshalError != nil {
// 			log.Println("unmarshal error", unmarshalError)
// 		}
// 		if thisPacket.Type == 100 {
// 			log.Println("100: User is already logged in.")
// 			conn.WriteMessage(1, []byte("You are already logged in."))
// 		} else if thisPacket.Type == 200 {
// 			log.Println("200: this is a friend request")
// 			conn.WriteMessage(1, []byte("friend request"))
// 		} else if thisPacket.Type == 300 {
// 			log.Println("300: this is a normal chat message")
// 			conn.WriteMessage(1, []byte("chat message"))
// 		} else {
// 			log.Println("???: I didn't understand this message")
// 			conn.WriteMessage(1, []byte("didn't understand this request"))
// 		}
// 	}

// }

func handleSuccessfulLogin(userpass []string, conn *websocket.Conn, msg *Packet) {
	username := userpass[0] + "@" + userpass[1]
	thisSession := UserSession{
		Username:   username,
		State:      ClientStateExchange{},
		Conn:       conn,
		Authorized: true,
	}
	if !fileExists("admin/users/" + userpass[0] + ".state") {
		createFile("admin/users/" + userpass[0] + ".state")
		thisSession.State = ClientStateExchange{
			CurrentFriends: []string{
				"fred@server.3ck0.com",
			},
			PendingFriends: []string{
				"mark@white.monster",
			},
			BlockedFriends: []string{},
			BlockedServers: []string{},
		}
		for _, r := range thisSession.State.CurrentFriends {
			for _, rr := range GlobalUserSessions {
				if rr.Username == r {
					r = "(" + r + ")"
				}
			}
		}
		marshalState, msErr := json.Marshal(thisSession.State)
		if msErr != nil {
			conn.Close()
			return
		}
		writeFile("admin/users/"+userpass[0]+".state", string(marshalState))
		conn.WriteJSON(marshalState)
	} else if stringExistsInFile(msg.Content) {
		thisFile := readFile("admin/users/" + userpass[0] + ".state")
		unmarshErr := json.Unmarshal([]byte(thisFile), &thisSession.State)
		for _, r := range thisSession.State.CurrentFriends {
			for _, rr := range GlobalUserSessions {
				if rr.Username == r {
					r = "(" + r + ")"
				}
			}
		}
		if unmarshErr != nil {
			log.Println("Unmarshal error: " + unmarshErr.Error())
		}
		conn.WriteMessage(1, []byte(thisFile))
	}
	AddUserSession(&thisSession)
	for i, session := range GlobalUserSessions {
		fmt.Printf(i + " " + session.Username + " ")
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
