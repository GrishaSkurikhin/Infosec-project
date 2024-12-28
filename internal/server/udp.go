package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/GrishaSkurikhin/divan_bot/internal/crypto/xor"
	"github.com/GrishaSkurikhin/divan_bot/internal/message"
)

type RequestHandler func(message []byte) error

type UDPServer struct {
	batchSize    int
	Conn         *net.UDPConn
	Handler      RequestHandler
	clientsMap   map[string][]message.Message
	symmetricKey []byte
	mu           sync.Mutex
}

func New(port int, batchSize int, symmetricKey []byte, handler RequestHandler) (*UDPServer, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		return nil, fmt.Errorf("net.ListenUDP: %w", err)
	}

	return &UDPServer{
		batchSize:    batchSize,
		Conn:         conn,
		Handler:      handler,
		clientsMap:   make(map[string][]message.Message),
		symmetricKey: symmetricKey,
	}, nil
}

func (s *UDPServer) Start() {
	defer s.Conn.Close()

	buffer := make([]byte, s.batchSize)
	for {
		s.processPacket(buffer)
	}
}

func (s *UDPServer) processPacket(buffer []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, clientAddr, err := s.Conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("conn.ReadFromUDP: %v\n", err)
		return
	}
	decryptedData := xor.Decrypt(buffer, s.symmetricKey)

	var msg message.Message
	if _, err := msg.Read(decryptedData); err != nil {
		log.Printf("message.Read: %v\n", err)
		return
	}
	s.clientsMap[clientAddr.String()] = append(s.clientsMap[clientAddr.String()], msg)

	log.Printf("Progress for client %s: %.2f", clientAddr.String(), msg.Progress())
	log.Printf("CRC check: %v", msg.Correct())

	if msg.Last() {
		var data []byte
		for _, batch := range s.clientsMap[clientAddr.String()] {
			data = append(data, batch.Data()...)
		}
		if err := s.Handler(data); err != nil {
			log.Printf("handler: %v\n", err)
		}
		delete(s.clientsMap, clientAddr.String())
	}
}
