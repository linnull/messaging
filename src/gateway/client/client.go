package client

import (
	"context"

	proto "gateway/proto"
	"gateway/server"

	"github.com/micro/go-micro/client"
)

type GatewayClient struct {
	cli proto.GatewayClient
}

func NewGatewayClient() *GatewayClient {
	return &GatewayClient{proto.NewGatewayClient(server.ServiceName, client.NewClient())}
}

func (client *GatewayClient) Echo(ctx context.Context, msg string) (string, error) {
	rsp, err := client.cli.Echo(ctx, &proto.EchoRequest{Msg: msg})
	if err != nil {
		return "", err
	}
	return rsp.Msg, nil
}
