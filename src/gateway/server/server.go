package server

import (
	"context"

	proto "gateway/proto"
)

type GatewayHandler struct{}

func (handler *GatewayHandler) Echo(ctx context.Context, req *proto.EchoRequest, rsp *proto.EchoResponse) error {
	rsp.Msg = req.Msg
	return nil
}
