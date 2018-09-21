package client

import (
	"context"

	"lib/micro"
	proto "logic/login/proto"
	"logic/login/server"
)

type LoginClient struct {
	cli proto.LoginClient
}

func NewLoginClient() *LoginClient {
	return &LoginClient{proto.NewLoginClient(server.ServiceName, micro.NewClient())}
}

func (client *LoginClient) Echo(ctx context.Context, msg string) (string, error) {
	rsp, err := client.cli.Echo(ctx, &proto.EchoRequest{Msg: msg})
	if err != nil {
		return "", err
	}
	return rsp.Msg, nil
}
