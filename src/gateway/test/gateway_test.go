package test

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"net"
	"testing"
	"time"

	"gateway/client"
	"lib/packet"
)

func TestClient(t *testing.T) {
	cli := client.NewGatewayClient()
	ctx := context.Background()
	msg, err := cli.Echo(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(msg)
}

func TestConn(t *testing.T) {
	ip := "127.0.0.1"
	port := 20000
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
}

func TestPacket(t *testing.T) {
	ip := "127.0.0.1"
	port := 20000
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	pbMsg := "test test test "
	msg := &packet.RequestMessage{
		UUID: 10086,
		Sequence: 1,
		MethodId: 1,
		Length: uint32(len(pbMsg)),
		PbMsgBytes: []byte(pbMsg),
		CheckSum: adler32.Checksum([]byte(pbMsg)),
	}
	messageBytes := make([]byte, 20+msg.Length)
	binary.BigEndian.PutUint64(messageBytes[0:8], msg.UUID)
	binary.BigEndian.PutUint32(messageBytes[8:12], msg.MethodId)
	binary.BigEndian.PutUint32(messageBytes[12:16], msg.Sequence)
	binary.BigEndian.PutUint32(messageBytes[16:20], msg.Length)
	copy(messageBytes[20:20+msg.Length], msg.PbMsgBytes)
	msg.CheckSum = adler32.Checksum(messageBytes)
	msgBytes, err := msg.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	pkg := &packet.Packet{
		MagicNumber:  packet.PACKET_MAGIC_NUM,
		ProtocolType: packet.PROTOCOL_TYPE_REQUEST,
		Length:       uint32(len(msgBytes)),
		Payload:      msgBytes,
	}
	pkgBytes, err := pkg.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	conn.Write(pkgBytes)
}
