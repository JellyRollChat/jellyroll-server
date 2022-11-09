package main

import (
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

const (
	webPort         = 5270
	trackerCommPort = 5269
	serverCommPort  = 5268
	clientCommPort  = 5267

	servertld = "server.3ck0.com"

	pubKeyFilePath    = "keys/" + "public.key"
	privKeyFilePath   = "keys/" + "private.key"
	signedKeyFilePath = "keys/" + "signed.key"
	selfCertFilePath  = "keys/" + "self.cert"
)

var serverKeys *ED25519Keys

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
	trackerPingLimit = 60 * time.Second
)
var (
	nc = ""

	brightblack   = ""
	brightred     = ""
	brightgreen   = ""
	brightyellow  = ""
	brightpurple  = ""
	brightmagenta = ""
	brightcyan    = ""
	brightwhite   = ""

	black   = ""
	red     = ""
	green   = ""
	yellow  = ""
	purple  = ""
	magenta = ""
	cyan    = ""
	white   = ""
)

var (
	joinMsg         []byte = []byte("JOIN")
	ncasMsg         []byte = []byte("NCAS")
	capkMsg         []byte = []byte("CAPK")
	certMsg         []byte = []byte("CERT")
	pingMsg         []byte = []byte("PING")
	pongMsg         []byte = []byte("PONG")
	pubkMsg         []byte = []byte("PUBK")
	nsigMsg         []byte = []byte("NSIG")
	sendMsg         []byte = []byte("SEND")
	rtrnMsg         []byte = []byte("RTRN")
	numTx           int
	wantsClean      bool = false
	serverPubkey         = ""
	serverPrivkey        = ""
	serverSignedkey      = ""
	upgrader             = websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
)
