package main

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func hashit(element string) string {
	hashedData := []byte(element + "saltElement0" + "saltElement1" + "saltElement2")
	slot := make([]byte, 64)
	sha3.ShakeSum256(slot, hashedData)
	return hex.EncodeToString(slot[:])
}
