package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"github.com/GrishaSkurikhin/divan_bot/internal/client"
	"github.com/GrishaSkurikhin/divan_bot/internal/gui"
	customImage "github.com/GrishaSkurikhin/divan_bot/internal/image"
	"github.com/GrishaSkurikhin/divan_bot/internal/server"
)

const (
	serverPortEnv = "SERVER_PORT"
	clientPortEnv = "CLIENT_PORT"

	keyEnv  = "SYMMETRIC_KEY"
	keySize = 32

	batchSize = 1024
)

func main() {
	serverPort, clientPort, err := getPorts()
	if err != nil {
		log.Fatalf("getPort: %v", err)
	}
	symmetricKey := make([]byte, keySize)
	if _, err := base64.StdEncoding.Decode(symmetricKey, []byte(os.Getenv(keyEnv))); err != nil {
		log.Fatalf("base64.StdEncoding.Decode: %v", err)
	}

	clnt := client.New(clientPort, batchSize, symmetricKey)
	if err != nil {
		log.Fatalf("client.New: %v", err)
	}

	app := gui.New(
		gui.WithRandomImageGenerator(customImage.GenerateRandomImage),
		gui.WithSubmitImageHandler(clnt.SendImage),
	)

	srv, err := server.New(serverPort, batchSize, symmetricKey, func(data []byte) error {
		img, err := convertBytesToImage(data)
		if err != nil {
			return fmt.Errorf("convertBytesToImage: %w", err)
		}
		app.SetReceivedImage(img)
		return nil
	})
	if err != nil {
		log.Fatalf("server.New: %v", err)
	}

	go srv.Start()
	app.Run()
}

func getPorts() (int, int, error) {
	serverPort, err := strconv.Atoi(os.Getenv(serverPortEnv))
	if err != nil {
		return 0, 0, fmt.Errorf("strconv.Atoi: %w", err)
	}
	clientPort, err := strconv.Atoi(os.Getenv(clientPortEnv))
	if err != nil {
		return 0, 0, fmt.Errorf("strconv.Atoi: %w", err)
	}
	return serverPort, clientPort, nil
}

func convertBytesToImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("image.Decode: %w", err)
	}
	return img, nil
}
