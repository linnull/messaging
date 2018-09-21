package test

import (
	"context"
	"fmt"
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
	pkgBytes, err := packet.MarshalRequestPacket(10086, 1, 1, []byte(pbMsg))
	if err != nil {
		t.Fatal(err)
	}
	conn.Write(pkgBytes)
}
