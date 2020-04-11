package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// ReadPacket reads a single packet from the given reader.
func ReadPacket(reader io.Reader, addr net.Addr, state State) (*Packet, error) {
	if state == StateHandshaking {
		bufReader := bufio.NewReader(reader)
		data, err := bufReader.Peek(1)
		if err != nil {
			return nil, err
		}

		if data[0] == LegacyServerListPingID {
			return readLegacyServerListPing(bufReader, addr)
		}
		reader = bufReader
	}

	frame, err := readFrame(reader, addr)
	if err != nil {
		return nil, err
	}

	packet := &Packet{Length: frame.length}

	remainder := bytes.NewBuffer(frame.payload)

	packet.PacketID, err = readVarInt(remainder)
	if err != nil {
		return nil, err
	}
	packet.Data = remainder.Bytes()
	logrus.WithFields(logrus.Fields{
		"client": addr,
		"packet": packet,
	}).Trace("read packet")

	return packet, nil
}

func readLegacyServerListPing(reader *bufio.Reader, addr net.Addr) (*Packet, error) {
	packetID, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if packetID != LegacyServerListPingID {
		return nil, errors.Errorf("expected legacy server listing ping packet ID, got %x", packetID)
	}

	payload, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if payload != 0x01 {
		return nil, errors.Errorf("expected payload=1 from legacy server listing ping, got %x", payload)
	}

	packetIDForPluginMsg, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if packetIDForPluginMsg != 0xFA {
		return nil, errors.Errorf("expected packetIDForPluginMsg=0xFA from legacy server listing ping, got %x", packetIDForPluginMsg)
	}

	messageNameShortLen, err := readUnsignedShort(reader)
	if err != nil {
		return nil, err
	}
	if messageNameShortLen != 11 {
		return nil, errors.Errorf("expected messageNameShortLen=11 from legacy server listing ping, got %d", messageNameShortLen)
	}

	messageName, err := readUTF16BEString(reader, messageNameShortLen)
	if messageName != "MC|PingHost" {
		return nil, errors.Errorf("expected messageName=MC|PingHost, got %s", messageName)
	}

	remainingLen, err := readUnsignedShort(reader)
	remainingReader := io.LimitReader(reader, int64(remainingLen))

	protocolVersion, err := readByte(remainingReader)
	if err != nil {
		return nil, err
	}

	hostnameLen, err := readUnsignedShort(remainingReader)
	if err != nil {
		return nil, err
	}
	hostname, err := readUTF16BEString(remainingReader, hostnameLen)
	if err != nil {
		return nil, err
	}

	port, err := readUnsignedInt(remainingReader)
	if err != nil {
		return nil, err
	}

	return &Packet{
		PacketID: LegacyServerListPingID,
		Length:   0,
		Data: &LegacyServerListPing{
			ProtocolVersion: int(protocolVersion),
			ServerAddress:   hostname,
			ServerPort:      uint16(port),
		},
	}, nil
}

func readUTF16BEString(reader io.Reader, symbolLen uint16) (string, error) {
	bsUtf16be := make([]byte, symbolLen*2)

	_, err := io.ReadFull(reader, bsUtf16be)
	if err != nil {
		return "", err
	}

	result, _, err := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), bsUtf16be)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func readFrame(reader io.Reader, addr net.Addr) (*Frame, error) {
	var err error
	frame := &Frame{}

	frame.length, err = readVarInt(reader)
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"client": addr,
		"length": frame.length,
	}).Trace("read frame length")

	frame.payload = make([]byte, frame.length)
	total := 0
	for total < frame.length {
		readIntoThis := frame.payload[total:]
		n, err := reader.Read(readIntoThis)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
		total += n
		logrus.WithFields(logrus.Fields{
			"client": addr,
			"total":  total,
		}).Trace("read frame content")

		if n == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	logrus.WithFields(logrus.Fields{
		"client": addr,
		"frame":  frame,
	}).Trace("read frame")
	return frame, nil
}

func readVarInt(reader io.Reader) (int, error) {
	b := make([]byte, 1)
	var numRead uint = 0
	result := 0
	for numRead <= 5 {
		n, err := reader.Read(b)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}
		value := b[0] & 0x7F
		result |= int(value) << (7 * numRead)

		numRead++

		if b[0]&0x80 == 0 {
			return result, nil
		}
	}

	return 0, errors.New("VarInt is too big")
}

func readString(reader io.Reader) (string, error) {
	length, err := readVarInt(reader)
	if err != nil {
		return "", err
	}

	b := make([]byte, 1)
	var strBuilder strings.Builder
	for i := 0; i < length; i++ {
		n, err := reader.Read(b)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}
		strBuilder.WriteByte(b[0])
	}

	return strBuilder.String(), nil
}

func readByte(reader io.Reader) (byte, error) {
	buf := make([]byte, 1)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func readUnsignedShort(reader io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(reader, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func readUnsignedInt(reader io.Reader) (uint32, error) {
	var value uint32
	err := binary.Read(reader, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// ReadHandshake reads a Handshake packet from the given data.
func ReadHandshake(data interface{}) (*Handshake, error) {

	dataBytes, ok := data.([]byte)
	if !ok {
		return nil, errors.New("data is not expected byte slice")
	}

	handshake := &Handshake{}
	buffer := bytes.NewBuffer(dataBytes)
	var err error

	handshake.ProtocolVersion, err = readVarInt(buffer)
	if err != nil {
		return nil, err
	}

	handshake.ServerAddress, err = readString(buffer)
	if err != nil {
		return nil, err
	}

	handshake.ServerPort, err = readUnsignedShort(buffer)
	if err != nil {
		return nil, err
	}

	nextState, err := readVarInt(buffer)
	if err != nil {
		return nil, err
	}
	handshake.NextState = nextState
	return handshake, nil
}
