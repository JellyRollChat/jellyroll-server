package main

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

func (s *FedServer) SyncMessages() {
	for {
		select {
		case msg := <-s.Inbox:
			// handle incoming message
			s.Messages[msg.ID] = msg
			if msg.RecipientServerURL == s.URL {
				// message is for this server, send it to recipient's inbox
				recipientServer := GlobalFedServers[msg.RecipientServerURL]
				recipientServer.Inbox <- msg
			}
		case msg := <-s.Outbox:
			// send message to other servers
			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("json marshal error: ", err)
				continue
			}
			serverURL := msg.RecipientServerURL
			if serverURL == "" {
				// message is for all servers, send it to all servers
				for _, server := range GlobalFedServers {
					if server.URL != s.URL {
						// send the message over the network
						conn, err := websocket.Dial("ws://"+server.URL+"/messages", "", "http://"+server.URL)
						if err != nil {
							log.Println("websocket dial error: ", err)
							continue
						}
						if _, err := conn.Write(data); err != nil {
							log.Println("websocket write error: ", err)
							continue
						}
						conn.Close()
					}
				}
			} else {
				// message is for a specific server, send it to that server
				server := GlobalFedServers[serverURL]
				server.Outbox <- msg
			}
		}
	}
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
