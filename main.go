package main

import (
	"fmt"
	"strconv"
)

func main() {

	announce("Server up!")

	createDirIfItDontExist("keys")

	fmt.Println("Server Pubkey:\t" + serverKeys.publicKey[0:4] + "..  ./" + pubKeyFilePath)
	fmt.Println("Server Privkey:\t" + serverKeys.privateKey[0:4] + "..  ./" + privKeyFilePath)
	fmt.Println("Server Sigkey:\t" + serverKeys.signedKey[0:4] + "..  ./" + signedKeyFilePath)
	fmt.Println("Server Cert:\t" + serverKeys.selfCert[0:4] + "..  ./" + selfCertFilePath)

	// open tracker channel
	announce("Tracker Socket Up!")
	go TrackerSocketAPI(serverKeys)
	fmt.Println("Tracker Port:\t" + strconv.Itoa(trackerCommPort))

	// open server channel
	announce("Server Socket Up!")
	go TrackerSocketAPI(serverKeys)
	fmt.Println("Server Port:\t" + strconv.Itoa(serverCommPort))

	// open client channel
	announce("Client Socket Up!")
	go TrackerSocketAPI(serverKeys)
	fmt.Println("Client Port:\t" + strconv.Itoa(clientCommPort))

	// blocking operation
	select {}
}

func announce(message string) {
	fmt.Printf("\n" + nc + green + "╔══════════════════════════════════════════╗\n   " + brightcyan + "+ " + message + green + "\n╚══════════════════════════════════════════╝\n\n" + nc)
}
