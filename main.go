package main

import (
	"fmt"
	"strconv"
)

func main() {

	announce("Server up!")

	createDirIfItDontExist("keys")
	createDirIfItDontExist("admin")
	createDirIfItDontExist("admin/users")

	osCheck()
	serverKeys = initKeys()

	fmt.Println("Server Pubkey:\t" + serverKeys.publicKey[0:4] + "..  ./" + pubKeyFilePath)
	fmt.Println("Server Privkey:\t" + serverKeys.privateKey[0:4] + "..  ./" + privKeyFilePath)
	fmt.Println("Server Sigkey:\t" + serverKeys.signedKey[0:4] + "..  ./" + signedKeyFilePath)
	fmt.Println("Server Cert:\t" + serverKeys.selfCert[0:4] + "..  ./" + selfCertFilePath)

	announce("Sockets Up!")

	// open client channel
	go SocketAPI(serverKeys)
	fmt.Println("Socket Port:\t" + strconv.Itoa(clientCommPort))

	announce("Web Frontend Up!")
	go serverWebAPI()
	fmt.Println("Web Port:\t" + strconv.Itoa(webPort))

	// blocking operation
	select {}
}
