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
	pingMsg  []byte = []byte("PING")
	upgrader        = websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
)
