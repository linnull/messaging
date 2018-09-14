package client

import (
	"context"

	"lib/micro"
	proto "template/proto"
	"template/server"
)

type TemplateClient struct {
	cli proto.TemplateClient
}

func NewTemplateClient() *TemplateClient {
	return &TemplateClient{proto.NewTemplateClient(server.ServiceName, micro.NewClient())}
}

func (client *TemplateClient) Echo(ctx context.Context, msg string) (string, error) {
	rsp, err := client.cli.Echo(ctx, &proto.EchoRequest{Msg: msg})
	if err != nil {
		return "", err
	}
	return rsp.Msg, nil
}
