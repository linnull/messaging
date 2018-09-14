package client

import (
	"context"

	proto "gateway/proto"
	"gateway/server"
	"lib/micro"
)

type GatewayClient struct {
	cli proto.GatewayClient
}

func NewGatewayClient() *GatewayClient {
	return &GatewayClient{proto.NewGatewayClient(server.ServiceName, micro.NewClient())}
}

func (client *GatewayClient) Echo(ctx context.Context, msg string) (string, error) {
	rsp, err := client.cli.Echo(ctx, &proto.EchoRequest{Msg: msg})
	if err != nil {
		return "", err
	}
	return rsp.Msg, nil
}
