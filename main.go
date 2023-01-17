package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {

	startTime = time.Now()
	osCheck()

	announce("JellyRoll Server")

	createDirIfItDontExist("keys")
	createDirIfItDontExist("admin")
	createDirIfItDontExist("admin/users")

	serverKeys = initKeys()

	announce("Sockets Information")

	// client comms socket
	go SocketAPI(serverKeys)
	fmt.Println("Client Socket Port:\t" + strconv.Itoa(clientCommPort))

	// server federation comms socket
	go FederationAPI(serverKeys)
	fmt.Println("Server Socket Port:\t" + strconv.Itoa(fedCommPort))

	announce("Web API Information")

	// HTTP API
	go WebAPI()
	fmt.Println("Web Port:\t" + strconv.Itoa(webPort))

	// blocking operation
	select {}
}
