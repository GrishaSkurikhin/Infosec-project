package main

import (
	"log"

	"github.com/GrishaSkurikhin/divan_bot/internal/crypto/xor"
)

func main() {
	symmetricKey, err := xor.GenerateKey()
	if err != nil {
		log.Fatalf("xor.GenerateKey: %v", err)
	}
	log.Printf("Symmetric key: %s\n", symmetricKey)
}
