package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

func socketAuthAgent(conn *websocket.Conn, keyCollection *ED25519Keys) {

	for {

		defer conn.Close()

		msgType, msg, err := conn.ReadMessage()

		if err != nil {

			// socket session closed
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())

			return
		}

		// if the socket has a join message
		if bytes.HasPrefix(msg, joinMsg) {

			msgToString := string(msg)

			trimNewlinePrefix := strings.TrimRight(msgToString, "\n")

			trimmedPubKey := strings.TrimLeft(trimNewlinePrefix, "JOIN ")

			// complains about loop duration
			regValidate, _ := regexp.MatchString(`[a-f0-9]{64}`, trimmedPubKey[:64])
			if !regValidate {
				fmt.Printf("\nContains illegal characters")
				conn.Close()
				return
			}

		}
		if bytes.HasPrefix(msg, rtrnMsg) {
			fmt.Printf("\nreturn message: %s", string(msg))

			input := strings.TrimLeft(string(msg), "RTRN ")
			var cert = strings.Split(input, " ")

			if !verifySignature(cert[0], keyCollection.publicKey, cert[1]) {
				fmt.Printf("\nsig doesnt verify")
			}
			if verifySignature(cert[0], keyCollection.publicKey, cert[1]) {
				fmt.Printf("\nsig verifies")
			}
		}
		if bytes.HasPrefix(msg, pubkMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}
		if bytes.HasPrefix(msg, peerMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}

	}
}

func sessionAgent(conn *websocket.Conn, sessionPubKey string) {
	fmt.Printf(brightgreen+"\n[%s] [%s] New socket session"+white, timeStamp(), conn.RemoteAddr())
	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] socket: %s\n"+white, timeStamp(), conn.RemoteAddr(), err)
			break
		}
		fmt.Println(msg)
	}
}
