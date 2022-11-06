package main

// // Download a set of keys as keys.txt for new users
// // domain.tld:5270/signup
// func SignupHandler(w http.ResponseWriter, r *http.Request) {

// 	// Create a set of keys each time this endpoint is loaded
// 	theseKeys := webKeys()

// 	// Set content type header for text
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	// Content-disposition tells the browser to download a file rather
// 	// than print to the screen like a web page
// 	w.Header().Set("Content-Disposition", "attachment; filename=\"keys-"+theseKeys.publicKey[0:4]+".txt\"")

// 	// Every keys.txt is the same size, so we declare it
// 	w.Header().Set("Content-Length", "484")

// 	// Assemble this into the header
// 	w.WriteHeader(http.StatusOK)

// 	// Announce that someone has hit this endpoint
// 	reportRequest("keys/new", w, r)

// 	// This is the content of keys-NNNN.txt
// 	// We are appending the first 4 characters of the pubkey to the filename
// 	// so that repeatedly loading the page doesn't accidentally overwrite
// 	// keys that are already in use.
// 	newAcctInfo := "pubkey:" + theseKeys.publicKey + "\nprivkey:" + theseKeys.privateKey + "\nselfcert:" + theseKeys.selfCert + "\nsignedkey" + theseKeys.signedKey

// 	// Write the full response with header and serve to the user
// 	w.Write([]byte(newAcctInfo))
// }

// func webKeys() *ED25519Keys {
// 	keys := ED25519Keys{}
// 	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
// 	if err != nil {
// 		handle("error: ", err)
// 	}
// 	keys.privateKey = hex.EncodeToString(privKey[0:32])
// 	keys.publicKey = hex.EncodeToString(pubKey)
// 	signedKey := ed25519.Sign(privKey, pubKey)
// 	keys.signedKey = hex.EncodeToString(signedKey)
// 	keys.selfCert = keys.publicKey + keys.signedKey
// 	return &keys
// }
