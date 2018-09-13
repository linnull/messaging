package packet

import (
	"encoding/binary"
	"errors"
	"hash/adler32"
)

const MESSAGE_MAX_LENGTH = 0x500000 // 5M

var (
	ErrMessageLengthLimit            = errors.New("message length limit exceeded")
	ErrUnexpectedPbMessageLength   = errors.New("unexpected pb message length")
	ErrUnexpectedPbMessageChecksum = errors.New("unexpected pb message checksum")
)

type RequestMessage struct {
	UUID       uint64
	MethodId   uint32
	Sequence   uint32
	Length     uint32
	PbMsgBytes []byte
	CheckSum   uint32
}

func (msg *RequestMessage) Marshal() (messageBytes []byte, err error) {
	msgBytesLen := msg.Length + 24
	if msgBytesLen > MESSAGE_MAX_LENGTH {
		return nil, ErrMessageLengthLimit
	}
	messageBytes = make([]byte, msgBytesLen)
	binary.BigEndian.PutUint64(messageBytes[0:8], msg.UUID)
	binary.BigEndian.PutUint32(messageBytes[8:12], msg.MethodId)
	binary.BigEndian.PutUint32(messageBytes[12:16], msg.Sequence)
	binary.BigEndian.PutUint32(messageBytes[16:20], msg.Length)
	copy(messageBytes[20:20+msg.Length], msg.PbMsgBytes)
	binary.BigEndian.PutUint32(messageBytes[20+msg.Length:24+msg.Length], msg.CheckSum)
	return messageBytes, nil
}

type ResponseMessage struct {
	UUID       uint64
	MethodId   uint32
	Sequence   uint32
	RetCode    uint32
	Length     uint32
	PbMsgBytes []byte
	CheckSum   uint32
}

type CommandMessage struct {
	UUID       uint64
	MethodId   uint32
	Length     uint32
	PbMsgBytes []byte
	CheckSum   uint32
}

func GetRequestMessageFromPacket(pkg *Packet) (msg *RequestMessage, err error) {
	if pkg.ProtocolType != PROTOCOL_TYPE_REQUEST {
		return nil, ErrUnsupportedPacketProtoType
	}
	if pkg.Length != uint32(len(pkg.Payload)) || pkg.Length < 20 {
		return nil, ErrUnexpectedPacketPayloadLength
	}
	msg = new(RequestMessage)
	msg.UUID = binary.BigEndian.Uint64(pkg.Payload[0:8])
	msg.MethodId = binary.BigEndian.Uint32(pkg.Payload[8:12])
	msg.Sequence = binary.BigEndian.Uint32(pkg.Payload[12:16])
	msg.Length = binary.BigEndian.Uint32(pkg.Payload[16:20])
	if msg.Length+24 != pkg.Length {
		return nil, ErrUnexpectedPbMessageLength
	}
	msg.PbMsgBytes = pkg.Payload[20 : 20+msg.Length]
	msg.CheckSum = binary.BigEndian.Uint32(pkg.Payload[20+msg.Length : 24+msg.Length])
	if msg.CheckSum != adler32.Checksum(pkg.Payload[0:20+msg.Length]) {
		return nil, ErrUnexpectedPbMessageChecksum
	}
	return msg, nil
}

func GetResponseMessageFromPacket(pkg *Packet) (msg *ResponseMessage, err error) {
	if pkg.ProtocolType != PROTOCOL_TYPE_RESPONSE {
		return nil, ErrUnsupportedPacketProtoType
	}
	if pkg.Length != uint32(len(pkg.Payload)) || pkg.Length < 28 {
		return nil, ErrUnexpectedPacketPayloadLength
	}
	msg = new(ResponseMessage)
	msg.UUID = binary.BigEndian.Uint64(pkg.Payload[0:8])
	msg.MethodId = binary.BigEndian.Uint32(pkg.Payload[8:12])
	msg.Sequence = binary.BigEndian.Uint32(pkg.Payload[12:16])
	msg.RetCode = binary.BigEndian.Uint32(pkg.Payload[16:20])
	msg.Length = binary.BigEndian.Uint32(pkg.Payload[20:24])
	if msg.Length+28 != pkg.Length {
		return nil, ErrUnexpectedPbMessageLength
	}
	msg.PbMsgBytes = pkg.Payload[24 : 24+msg.Length]
	msg.CheckSum = binary.BigEndian.Uint32(pkg.Payload[24+msg.Length : 28+msg.Length])
	if msg.CheckSum != adler32.Checksum(pkg.Payload[0:24+msg.Length]) {
		return nil, ErrUnexpectedPbMessageChecksum
	}
	return msg, nil
}

func GetCommandMessageFromPacket(pkg *Packet) (msg *CommandMessage, err error) {
	if pkg.ProtocolType != PROTOCOL_TYPE_COMMAND {
		return nil, ErrUnsupportedPacketProtoType
	}
	if pkg.Length != uint32(len(pkg.Payload)) || pkg.Length < 20 {
		return nil, ErrUnexpectedPacketPayloadLength
	}
	msg = new(CommandMessage)
	msg.UUID = binary.BigEndian.Uint64(pkg.Payload[0:8])
	msg.MethodId = binary.BigEndian.Uint32(pkg.Payload[8:12])
	msg.Length = binary.BigEndian.Uint32(pkg.Payload[12:16])
	if msg.Length+20 != pkg.Length {
		return nil, ErrUnexpectedPbMessageLength
	}
	msg.PbMsgBytes = pkg.Payload[16 : 16+msg.Length]
	msg.CheckSum = binary.BigEndian.Uint32(pkg.Payload[16+msg.Length : 20+msg.Length])
	if msg.CheckSum != adler32.Checksum(pkg.Payload[0:16+msg.Length]) {
		return nil, ErrUnexpectedPbMessageChecksum
	}
	return msg, nil
}
