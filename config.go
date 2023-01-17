package main

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	webPort        = 5270
	clientCommPort = 5267
	fedCommPort    = 5264

	servertld = "server.3ck0.com"

	pubKeyFilePath    = "keys/" + "public.key"
	privKeyFilePath   = "keys/" + "private.key"
	signedKeyFilePath = "keys/" + "signed.key"
	selfCertFilePath  = "keys/" + "self.cert"
)

var (
	startTime   time.Time
	userCount   int
	socketCount int
)

var serverKeys *ED25519Keys

var GlobalUserSessions = make(map[string]*UserSession)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
)

// Packet is an object to encapsulate messages
type Packet struct {
	Type    int    `json:"msg_type"`
	Content string `json:"msg_content"`
}

// AuthObject struct represents an object used for authentication purposes. It has two fields: a string Username representing the username of the user, and a string Password representing the password of the user.
type AuthObject struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MessageHistoryEntry struct represents an entry in the message history for a user. It has two fields: an int64 Timestamp representing the time the message was sent or received, and a string Message representing the content of the message.
type MessageHistoryEntry struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// UserSession struct represents a session for a user in a messaging system. It has seven fields:
// Username: a string representing the username of the user
// State: a ClientStateExchange object representing the current state of the user's client
// Conn: a pointer to a websocket connection representing the connection to the user's client
// Authorized: a bool representing whether the user is authorized to use the service
// History: a slice of MessageHistoryEntry structs representing the message history for the user
// Inbox: a channel for incoming messages for the user
// Outbox: a channel for outgoing messages from the user
type UserSession struct {
	Username   string                `json:"username"`
	State      ClientStateExchange   `json:"state"`
	Conn       *websocket.Conn       `json:"conn"`
	Authorized bool                  `json:"authorized"`
	History    []MessageHistoryEntry `json:"history"`
	Inbox      chan *FedMessage      `json:"inbox"`
	Outbox     chan *FedMessage      `json:"outbox"`
}

var mutex sync.Mutex

// ClientMessage is a simple format for basic user<->user messages that are passed through a server
type ClientMessage struct {
	Timestamp int64  `json:"timestamp"`
	From      string `json:"from"` // alice@server1.tld sending the message
	Recv      string `json:"recv"` // bob@server2.tld receiving the message
	Body      string `json:"body"` // the message body
}

// ClientStateExchange is an interaction with the server that conveys busy status, current friends list, unconfirmed friend requests, blocked users and blocked servers. When a friend request is received from the server, that friend ID is added to the PendingFriends. If it is accepted, the friend ID is added to Friends and removed from PendingFriends, then a ClientStateExchange is sent back to the server to reflect the change. Rejected friend request does not add to BlockedFriends, but the user is presented with accept, reject, block menu.
// CurrentFriends: json list of current friends
// PendingFriends: unconfirmed, denied friend requests
// BlockedFriends: drop messages from these users
// BlockedServers: drop messages from these servers
type ClientStateExchange struct {
	PendingFriends []string `json:"pending_friends"`
	CurrentFriends []string `json:"current_friends"`
	BlockedFriends []string `json:"blocked_friends"`
	BlockedServers []string `json:"blocked_servers"`
}

// FedMessage struct includes fields for the message ID, the sender and recipient IDs, the message content, a timestamp for when the message was sent, and the URLs for the sender and recipient servers. The fields for the sender and recipient servers are necessary in a federated messaging system to keep track of where the message should be sent and where it came from.
type FedMessage struct {
	ID                 string    `json:"id"`
	SenderID           string    `json:"sender_id"`
	RecipientID        string    `json:"recipient_id"`
	Content            string    `json:"content"`
	Timestamp          time.Time `json:"timestamp"`
	SenderServerURL    string    `json:"sender_server_url"`
	RecipientServerURL string    `json:"recipient_server_url"`
}

// GlobalFedServers is a map of FedServer objects, indexed by the URL of the server.
var GlobalFedServers = make(map[string]*FedServer)

// FedServer represents a federated server in the messaging system. It has four fields:
//   - URL is the URL of the server.
//   - Inbox is a channel for incoming messages for the server.
//   - Outbox is a channel for outgoing messages from the server.
//   - Messages is a map of FedMessage objects indexed by their ID.
//   - Websocket is a pointer to a websocket connection for the server.
type FedServer struct {
	URL       string                 `json:"url"`
	Inbox     chan *FedMessage       `json:"inbox"`
	Outbox    chan *FedMessage       `json:"outbox"`
	Messages  map[string]*FedMessage `json:"messages"`
	Websocket *websocket.Conn        `json:"websocket"`
}

// FriendRequest struct represents a request to add a user as a friend. It has two fields: a string From representing the username of the sender, and a string To representing the username of the recipient.
type FriendRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}
