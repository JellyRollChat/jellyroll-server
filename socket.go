package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type CommPacket struct {
}

func SocketAPI(keyCollection *ED25519Keys) {

	api := mux.NewRouter()

	// Channel Socket
	api.HandleFunc("/talk", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] [%s] +client\n"+white, timeStamp(), conn.RemoteAddr())
		socketParser(conn, keyCollection)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(clientCommPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func socketParser(conn *websocket.Conn, keyCollection *ED25519Keys) {

	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// socket session closed
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
			break
		}

		msgType := 1

		if bytes.HasPrefix(msg, pingMsg) {
			conn.WriteMessage(msgType, []byte(pingMsg))
		} else if bytes.HasPrefix(msg, infoMsg) {
			fmt.Println("Full User ID:")
			fullnameb := hex.EncodeToString([]byte("username@server.tld"))
			fmt.Println(fullnameb)
			fmt.Println("Short nickname: ")
			shortnameb := hex.EncodeToString([]byte("username"))
			fmt.Println(shortnameb)
			fmt.Println("Body: ")
			bodyb := hex.EncodeToString([]byte("My public key is: " + keyCollection.publicKey))
			fmt.Println(bodyb)
			fmt.Println("Full Message: ")
			fullmesgb := "<:" + fullnameb + "," + bodyb + ":>"
			fmt.Println(fullmesgb)
			conn.WriteMessage(msgType, []byte(fullmesgb))
		} else if bytes.HasPrefix(msg, rglrMsg) {

			msgStr := string(msg)

			trimLCarrot := strings.TrimLeft(msgStr, "<:")
			trimRCarrot := strings.TrimRight(trimLCarrot, ":>")
			splitUsrMsg := strings.Split(trimRCarrot, ",")

			fmt.Println("Full User ID:")
			usernameEnc := fmt.Sprintf("%s", splitUsrMsg[0])
			fmt.Println("Encoded: ", usernameEnc)
			decodedUserName, decodeErr := hex.DecodeString(usernameEnc)
			if decodeErr != nil {
				fmt.Println("There was an error decoding the username. ")
			} else {
				fmt.Println("Decoded bytes: ", decodedUserName)
				fmt.Println("Decoded string: ", string(decodedUserName))
			}

			fmt.Println("Message Body:")
			bodyEnc := fmt.Sprintf("%s", splitUsrMsg[1])
			fmt.Println("Encoded: ", bodyEnc)
			decodedBody, decodeErr := hex.DecodeString(bodyEnc)
			if decodeErr != nil {
				fmt.Println("There was an error decoding the body. ")
			} else {
				fmt.Println("Decoded bytes: ", decodedBody)
				fmt.Println("Decoded string: ", string(decodedBody))
			}

			thisDecodedMesg := "<:" + string(decodedUserName) + "," + string(decodedBody) + ":>"

			conn.WriteMessage(msgType, []byte(thisDecodedMesg))
		}

	}
}
