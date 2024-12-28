package xor

import (
	"crypto/rand"
	"encoding/base64"
)

const keySize = 32

func GenerateKey() (string, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)
	return encodedKey, nil
}

func Encrypt(data, key []byte) []byte {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)]
	}
	return result
}

func Decrypt(encryptedData, key []byte) []byte {
	result := make([]byte, len(encryptedData))
	for i := 0; i < len(encryptedData); i++ {
		result[i] = encryptedData[i] ^ key[i%len(key)]
	}
	return result
}
