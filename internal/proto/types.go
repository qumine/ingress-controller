package proto

import "fmt"

// Frame represents a single frame in the communication.
type Frame struct {
	length  int
	payload []byte
}

// State represents the state a minecraft connection is in.
type State int

const (
	// StateHandshaking is the initial state of a minecraft connection.
	StateHandshaking = iota
)

var trimLimit = 64

func trimBytes(data []byte) ([]byte, string) {
	if len(data) < trimLimit {
		return data, ""
	}
	return data[:trimLimit], "..."
}

func (f *Frame) String() string {
	trimmed, cont := trimBytes(f.payload)
	return fmt.Sprintf("Frame:[len=%d, payload=%#X%s]", f.length, trimmed, cont)
}

// Packet represents a single packet in the minecraft communication.
type Packet struct {
	// Length is the length of the packet.
	Length int
	// PacketID is the ID of the packet.
	PacketID int
	// Data is either a byte slice of raw content or a parsed message
	Data interface{}
}

func (p *Packet) String() string {
	if dataBytes, ok := p.Data.([]byte); ok {
		trimmed, cont := trimBytes(dataBytes)
		return fmt.Sprintf("Frame:[len=%d, packetId=%d, data=%#X%s]", p.Length, p.PacketID, trimmed, cont)
	}
	return fmt.Sprintf("Frame:[len=%d, packetId=%d, data=%+v]", p.Length, p.PacketID, p.Data)

}

const (
	// HandshakeID is the ID of the Handshake packet.
	HandshakeID = 0x00
	// LegacyServerListPingID is the ID of the LegacyServerListPing packet.
	LegacyServerListPingID = 0xFE
)

// Handshake is the first packet in the minecraft protocol send by the client.
type Handshake struct {
	ProtocolVersion int
	ServerAddress   string
	ServerPort      uint16
	NextState       int
}

// LegacyServerListPing is send by legacy minecraft client.
type LegacyServerListPing struct {
	ProtocolVersion int
	ServerAddress   string
	ServerPort      uint16
}

type byteReader interface {
	ReadByte() (byte, error)
}
