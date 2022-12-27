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

	announce("Socket Up!")

	// open client channel
	go SocketAPI(serverKeys)
	fmt.Println("Socket Port:\t" + strconv.Itoa(clientCommPort))

	announce("Web Frontend Up!")
	go WebAPI()
	fmt.Println("Web Port:\t" + strconv.Itoa(webPort))

	// blocking operation
	select {}
}
