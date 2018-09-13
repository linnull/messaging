package packet

import (
	"encoding/binary"
	"errors"
	"io"
	"runtime"
)

const (
	PROTOCOL_TYPE_REQUEST  = 1
	PROTOCOL_TYPE_RESPONSE = 2
	PROTOCOL_TYPE_COMMAND  = 3

	PACKET_MAGIC_NUM  = 0x54FB
	PACKET_MAX_LENGTH = 0x500000 // 5M
)

var (
	ErrPacketLengthLimit            = errors.New("packet length limit exceeded")
	ErrUnexpectedPacketMagicNum     = errors.New("unexpected packet magic num")
	ErrUnsupportedPacketProtoType = errors.New("unsupported packet protocol type")
	ErrUnexpectedPacketPayloadLength = errors.New("unexpected packet payload length")
)

type Packet struct {
	MagicNumber  uint16
	ProtocolType uint16
	Length       uint32
	Payload      []byte
}

func (p *Packet) Marshal() (packetBytes []byte, err error) {
	msgBytesLen := uint32(len(p.Payload))
	packetBytesLen := msgBytesLen + 8
	if packetBytesLen > PACKET_MAX_LENGTH {
		return nil, ErrPacketLengthLimit
	}
	packetBytes = make([]byte, packetBytesLen)
	binary.BigEndian.PutUint16(packetBytes[0:2], p.MagicNumber)
	binary.BigEndian.PutUint16(packetBytes[2:4], p.ProtocolType)
	binary.BigEndian.PutUint32(packetBytes[4:8], msgBytesLen)
	copy(packetBytes[8:packetBytesLen], p.Payload)
	return packetBytes, nil
}

func GetPacket(reader io.Reader) (packet *Packet, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	pkgHeader := make([]byte, 8)
	_, err = io.ReadFull(reader, pkgHeader)
	if err != nil {
		return nil, err
	}
	magicNumber := binary.BigEndian.Uint16(pkgHeader[:2])
	if magicNumber != PACKET_MAGIC_NUM {
		return nil, ErrUnexpectedPacketMagicNum
	}

	protoType := binary.BigEndian.Uint16(pkgHeader[2:4])
	length := binary.BigEndian.Uint32(pkgHeader[4:8])
	if length > PACKET_MAX_LENGTH {
		return nil, ErrPacketLengthLimit
	}

	pkg := new(Packet)
	pkg.MagicNumber = magicNumber
	pkg.ProtocolType = protoType
	pkg.Length = length
	pkg.Payload = make([]byte, length)
	_, err = io.ReadFull(reader, pkg.Payload)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
