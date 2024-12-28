package client

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net"

	"github.com/GrishaSkurikhin/divan_bot/internal/crypto/xor"
	"github.com/GrishaSkurikhin/divan_bot/internal/message"
)

type UDPClient struct {
	port         int
	bufferSize   int
	symmetricKey []byte
}

func New(port int, bufferSize int, symmetricKey []byte) *UDPClient {
	return &UDPClient{
		port:         port,
		bufferSize:   bufferSize,
		symmetricKey: symmetricKey,
	}
}

func (c *UDPClient) SendImage(img image.Image) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 50}); err != nil {
		return fmt.Errorf("jpeg.Encode: %w", err)
	}
	if err := c.SendData(buf.Bytes()); err != nil {
		return fmt.Errorf("client.SendData: %w", err)
	}
	return nil
}

func (c *UDPClient) SendData(data []byte) error {
	batches, err := message.Batches(data, c.bufferSize)
	if err != nil {
		return fmt.Errorf("message.Batches: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: c.port})
	if err != nil {
		return fmt.Errorf("net.DialUDP: %w", err)
	}
	defer conn.Close()

	for _, batch := range batches {
		batchBytes := make([]byte, c.bufferSize)
		if _, err := batch.Write(batchBytes); err != nil {
			log.Printf("Error write batch: %v\n", err)
		}
		encryptedMessage := xor.Encrypt(batchBytes, c.symmetricKey)

		if _, err = conn.Write(encryptedMessage); err != nil {
			log.Printf("Error sending packet: %v\n", err)
		}
	}
	return nil
}
