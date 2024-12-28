package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/GrishaSkurikhin/divan_bot/internal/crc"
)

const (
	uint16Size = 2
	uint32Size = 4
)

type Message struct {
	crc            uint32
	totalSegments  uint16
	currentSegment uint16
	data           []byte
}

// significantPart возвращает значимую часть сообщения в байтах
func (m *Message) significantPart() []byte {
	data := make([]byte, len(m.data)+2*uint16Size)

	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+uint16Size], m.totalSegments)
	offset += uint16Size
	binary.BigEndian.PutUint16(data[offset:offset+uint16Size], m.currentSegment)
	offset += uint16Size

	copy(data[offset:], m.data)
	return data
}

func (m *Message) setCRC() {
	data := m.significantPart()
	m.crc = crc.ComputeCRC32(data)
}

func (m *Message) Correct() bool {
	return crc.ComputeCRC32(m.significantPart()) == m.crc
}

func Batches(data []byte, batchSize int) ([]Message, error) {
	var messages []Message

	metadataSize := uint16Size*2 + uint32Size
	maxDataSize := batchSize - metadataSize
	if maxDataSize <= 0 {
		return nil, errors.New("batchSize is too small")
	}

	totalSegments := (len(data) + maxDataSize - 1) / maxDataSize
	if totalSegments > math.MaxUint16 {
		return nil, fmt.Errorf("totalSegments %d слишком велико для типа uint16", totalSegments)
	}

	for i := 0; i < totalSegments; i++ {
		start := i * maxDataSize
		end := start + maxDataSize
		if end > len(data) {
			end = len(data)
		}

		messageData := make([]byte, maxDataSize)
		copy(messageData, data[start:end])

		message := Message{
			totalSegments:  uint16(totalSegments),
			currentSegment: uint16(i + 1),
			data:           messageData,
		}
		message.setCRC()
		messages = append(messages, message)
	}
	return messages, nil
}

// Write записывает данные объекта Message в переданный срез байтов
func (m *Message) Write(p []byte) (n int, err error) {
	minSize := 2*uint16Size + uint32Size
	if len(p) < minSize {
		return 0, fmt.Errorf("buffer too small: need at least %d bytes", minSize)
	}

	offset := 0
	binary.BigEndian.PutUint32(p[offset:offset+uint32Size], m.crc)
	offset += uint32Size

	binary.BigEndian.PutUint16(p[offset:offset+uint16Size], m.totalSegments)
	offset += uint16Size
	binary.BigEndian.PutUint16(p[offset:offset+uint16Size], m.currentSegment)
	offset += uint16Size

	copy(p[offset:], m.data)
	offset += len(m.data)

	return offset, nil
}

// Read восстанавливает данные объекта Message из переданного среза байтов
func (m *Message) Read(p []byte) (n int, err error) {
	minSize := 2*uint16Size + uint32Size
	if len(p) < minSize {
		return 0, fmt.Errorf("buffer too small: need at least %d bytes", minSize)
	}

	offset := 0
	m.crc = binary.BigEndian.Uint32(p[offset : offset+uint32Size])
	offset += uint32Size

	m.totalSegments = binary.BigEndian.Uint16(p[offset : offset+uint16Size])
	offset += uint16Size
	m.currentSegment = binary.BigEndian.Uint16(p[offset : offset+uint16Size])
	offset += uint16Size

	m.data = make([]byte, len(p)-offset)
	copy(m.data, p[offset:])
	offset += len(m.data)

	return offset, nil
}

func (m *Message) Data() []byte {
	return m.data
}

func (m *Message) Last() bool {
	return m.totalSegments == m.currentSegment
}

func (m *Message) Progress() float64 {
	if m.totalSegments == 0 {
		return math.NaN()
	}
	return float64(m.currentSegment) / float64(m.totalSegments) * 100.0
}
