package main

import "github.com/gorilla/websocket"

const (
	trackerCommPort = 5269
	serverCommPort  = 5268
	clientCommPort  = 5267

	pubKeyFilePath    = "keys/" + "public.key"
	privKeyFilePath   = "keys/" + "private.key"
	signedKeyFilePath = "keys/" + "signed.key"
	selfCertFilePath  = "keys/" + "self.cert"
)

var serverKeys *ED25519Keys

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
	dbUser          string = "postgres"
	dbName          string = "karai"
	dbSSL           string = "disable"
	joinMsg         []byte = []byte("JOIN")
	ncasMsg         []byte = []byte("NCAS")
	capkMsg         []byte = []byte("CAPK")
	certMsg         []byte = []byte("CERT")
	peerMsg         []byte = []byte("PEER")
	pubkMsg         []byte = []byte("PUBK")
	nsigMsg         []byte = []byte("NSIG")
	sendMsg         []byte = []byte("SEND")
	rtrnMsg         []byte = []byte("RTRN")
	numTx           int
	wantsClean      bool = false
	karaiAPIPort    int
	serverPubkey    = ""
	serverPrivkey   = ""
	serverSignedkey = ""
	upgrader        = websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
)
