package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// SyncMessages first checks if the Websocket field is nil, in which case it calls the InitWebsocket method to initialize the socket connection. If the Websocket field is not nil, it uses the WriteJSON function from the gorilla/websocket package to send the Messages field of the FedServer struct as JSON data to the remote server. It returns true if the sync is successful, and false if there is any error.
func (s *FedServer) SyncMessages() bool {
	// using the Websocket field of the FedServer struct, check for an active socket, or use InitWebsocket() method to create one if none exists.
	if s.Websocket == nil {
		if err := s.InitWebsocket(); err != nil {
			return false
		}
	}

	// use s.Websocket and json to sync all message channels with the remote FedServer
	if err := s.Websocket.WriteJSON(s.Messages); err != nil {
		return false
	}

	// return true if successful in syncing
	return true
}

// AddFedServer function adds a new server to the GlobalFedServers map. It creates a new FedServer struct with the provided serverURL and assigns it to the server variable. It then adds the server to the GlobalFedServers map with the server's URL as the key. The Inbox and Outbox channels and the Messages map are also initialized in this function.
func (s *FedServer) AddFedServer(serverURL string) {
	// add a new server to the map
	server := &FedServer{
		URL:      serverURL,
		Inbox:    make(chan *FedMessage),
		Outbox:   make(chan *FedMessage),
		Messages: make(map[string]*FedMessage),
	}
	GlobalFedServers[server.URL] = server
}

// GetFedServer function retrieves the server from the GlobalFedServers map with the provided serverURL as the key and returns it.
func (s *FedServer) GetFedServer(serverURL string) *FedServer {
	return GlobalFedServers[serverURL]
}

// SendMessageToFedUser is a function that takes a ClientMessage and sends it to a federated server over a socket using gorilla/websocket
func (s *FedServer) SendMessageToFedUser(msg ClientMessage) bool {
	if s.Websocket == nil {
		if err := s.InitWebsocket(); err != nil {
			log.Println("error initializing websocket:", err)
			return false
		}
	}
	if err := s.Websocket.WriteJSON(msg); err != nil {
		log.Println("error sending message:", err)
		return false
	}
	return true
}

func (s *FedServer) InitWebsocket() error {
	// Dial the server's websocket endpoint using the URL stored in the FedServer struct
	conn, _, err := websocket.DefaultDialer.Dial(s.URL, nil)
	if err != nil {
		return err
	}
	// Store the websocket connection in the FedServer struct
	s.Websocket = conn
	// Update the FedServer in the GlobalFedServers map
	GlobalFedServers[s.URL] = s
	return nil
}
