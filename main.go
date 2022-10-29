package main

import (
	"fmt"
	"strconv"
)

func main() {

	announce("Server up!")

	createDirIfItDontExist("keys")

	// serverKeys := initKeys()

	fmt.Println("Server Pubkey:\t" + serverKeys.publicKey[0:4] + "..  ./" + pubKeyFilePath)
	fmt.Println("Server Privkey:\t" + serverKeys.privateKey[0:4] + "..  ./" + privKeyFilePath)
	fmt.Println("Server Sigkey:\t" + serverKeys.signedKey[0:4] + "..  ./" + signedKeyFilePath)
	fmt.Println("Server Cert:\t" + serverKeys.selfCert[0:4] + "..  ./" + selfCertFilePath)

	announce("Socket API up!")

	go SocketAPI(serverKeys)

	fmt.Println("Tracker Port:\t" + strconv.Itoa(trackerCommPort))
	fmt.Println("Server Port:\t" + strconv.Itoa(serverCommPort))
	fmt.Println("Client Port:\t" + strconv.Itoa(clientCommPort))

	select {}
}

func announce(message string) {
	fmt.Printf("\n" + nc + green + "╔══════════════════════════════════════════╗\n   " + brightcyan + "+ " + message + green + "\n╚══════════════════════════════════════════╝\n\n" + nc)
}
