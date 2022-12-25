package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

const (
	webPort         = 5270
	trackerCommPort = 5269
	serverCommPort  = 5268
	clientCommPort  = 5267

	servertld     = "server.3ck0.com"
	defaultUser   = "admin"
	defaultSender = defaultUser + "@" + servertld

	pubKeyFilePath    = "keys/" + "public.key"
	privKeyFilePath   = "keys/" + "private.key"
	signedKeyFilePath = "keys/" + "signed.key"
	selfCertFilePath  = "keys/" + "self.cert"
)

var serverKeys *ED25519Keys

var (
	UserSessions   = make(chan UserSession)
	OutboxMessages = make(chan Packet)
)

var (
	corsAllowedHeaders = []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Cache-Control",
		"Content-Security-Policy",
		"Feature-Policy",
		"Referrer-Policy",
		"X-Requested-With"}

	corsOrigins = []string{
		// "*",
		"127.0.0.1"}

	corsMethods = []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"OPTIONS"}

	headersCORS = handlers.AllowedHeaders(corsAllowedHeaders)
	originsCORS = handlers.AllowedOrigins(corsOrigins)
	methodsCORS = handlers.AllowedMethods(corsMethods)
)

var (
	nc = ""

	// brightblack   = ""
	brightred    = ""
	brightgreen  = ""
	brightyellow = ""
	// brightpurple  = ""
	brightmagenta = ""
	brightcyan    = ""
	// brightwhite   = ""

	// black   = ""
	// red     = ""
	green = ""
	// yellow  = ""
	purple  = ""
	magenta = ""
	cyan    = ""
	white   = ""
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
)

// Packet is an object to encapsulate messages
// Types: 100 friend request, 200 normal message user to user
type Packet struct {
	MsgType    int    `json:"msg_type"`
	MsgContent string `json:"msg_content"`
}

type AuthObject struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSession struct {
	Username  string
	Servertld string
	Conn      *websocket.Conn
}

// Message is a simple format for basic user<->user messages that are passed through a server
type Message struct {
	From string `json:"from"`
	Recv string `json:"recv"`
	Body string `json:"body"`
}

// StateExchange is an interaction with the server that conveys busy status, current friends list, unconfirmed friend requests, blocked users and blocked servers. When a friend request is received from the server, that friend ID is added to the PendingFriends. If it is accepted, the friend ID is added to Friends and removed from PendingFriends, then a StateExchange is sent back to the server to reflect the change. Rejected friend request does not add to BlockedFriends, but the user is presented with accept, reject, block menu.
// CurrentFriends: json list of current friends
// PendingFriends: unconfirmed, denied friend requests
// BlockedFriends: drop messages from these users
// BlockedServers: drop messages from these servers
type StateExchange struct {
	PendingFriends []string `json:"pending_friends"`
	CurrentFriends []string `json:"current_friends"`
	BlockedFriends []string `json:"blocked_friends"`
	BlockedServers []string `json:"blocked_servers"`
}
