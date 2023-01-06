package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// FederationAPI is a function that sets up an HTTP server to listen for websocket connections on a specific port, with CORS headers set to allow all origins and methods. When a connection is received, it upgrades the connection to a websocket and passes it to the fedSocketHandler function.
func FederationAPI(keyCollection *ED25519Keys) {
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
	api.HandleFunc("/federate", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] +server [%s]\n"+nc, timeStamp(), conn.RemoteAddr())
		fedSocketHandler(conn, keyCollection)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(fedCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

// fedSocketHandler is a function that is called when a new websocket connection is received on the HTTP server setup by the FederationAPI function. It creates a new FedServer object with the connection's URL, a new Inbox and Outbox channel, an empty Messages map, and the websocket connection. It then adds the FedServer object to the GlobalFedServers map. It then enters a loop to listen for incoming messages on the websocket connection. If the message is intended for this server, it is added to the FedServer's Inbox channel. If the message is not intended for this server, it is forwarded to the appropriate server by sending it through the Outbox channel of the destination FedServer object in the GlobalFedServers map.
func fedSocketHandler(conn *websocket.Conn, keyCollection *ED25519Keys) {
	// Create a new FedServer object
	f := &FedServer{
		URL:       conn.RemoteAddr().String(),
		Inbox:     make(chan *FedMessage),
		Outbox:    make(chan *FedMessage),
		Messages:  make(map[string]*FedMessage),
		Websocket: conn,
	}

	// Add the FedServer object to the GlobalFedServers map
	GlobalFedServers[f.URL] = f

	// Start listening for messages on the websocket connection
	for {
		var msg FedMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("error reading message:", err)
			break
		}

		// Check if the message is intended for this server
		if msg.RecipientServerURL == f.URL {
			// If the message is intended for this server, add it to the FedServer's inbox
			f.Inbox <- &msg
		} else {
			// If the message is not intended for this server, forward it to the appropriate server
			server, ok := GlobalFedServers[msg.RecipientServerURL]
			if !ok {
				log.Println("error: recipient's server is not registered")
				continue
			}
			server.Outbox <- &msg
		}
	}
}

// Listen is a method of the FedServer struct that listens for incoming messages on the Inbox and Outbox channels and processes them accordingly. If a message is received on the Inbox channel, it is processed and sent to the appropriate user. If a message is received on the Outbox channel, it is forwarded to the appropriate server.
func (f *FedServer) Listen() {
	for {
		select {
		case msg := <-f.Inbox:
			// Check if the message is intended for this server
			if msg.RecipientServerURL != servertld {
				log.Println("error: message intended for different server")
				continue
			}

			// Check if the user has an active session
			userSession, ok := GlobalUserSessions[msg.RecipientID]
			if !ok {
				// If the user doesn't have an active session, check if they exist in users.list
				if ok := stringExistsInFile(msg.RecipientID); !ok {
					log.Println("error: recipient user not found")
					continue
				}

				// If the user exists in users.list, send a message indicating that the user is offline
				err := f.Websocket.WriteJSON(&FedMessage{
					ID:                 "", // ID field can be left empty
					SenderID:           servertld,
					RecipientID:        msg.SenderID,
					Timestamp:          time.Now(), // You can set the Timestamp field to the current time
					SenderServerURL:    "",         // SenderServerURL field can be left empty
					RecipientServerURL: "",         // RecipientServerURL field can be left empty
					Content:            "",
				})

				if err != nil {
					log.Println("error sending message to sender:", err)
					continue
				}
				continue
			}

			// If the user has an active session, add the message to their Inbox channel
			userSession.Inbox <- msg
		}
	}
}

// Send is a method of the FedServer struct that listens for outgoing messages on the Outbox channel and sends them to the appropriate server. If a message is received on the Outbox channel, it is forwarded to the appropriate server. If there is an error sending the message, it is logged. This method runs indefinitely until the FedServer is shutdown.
func (f *FedServer) Send() {
	for {
		select {
		case msg := <-f.Outbox:
			// Send the message to the recipient server
			server, ok := GlobalFedServers[msg.RecipientServerURL]
			if !ok {
				log.Println("error: recipient server not found")
				continue
			}
			err := server.Websocket.WriteJSON(msg)
			if err != nil {
				log.Println("error sending message to recipient server:", err)
				continue
			}
		}
	}
}
