package server

import (
	"context"

	proto "template/proto"
)

type TemplateHandler struct{}

func (handler *TemplateHandler) Echo(ctx context.Context, req *proto.EchoRequest, rsp *proto.EchoResponse) error {
	rsp.Msg = req.Msg
	return nil
}
