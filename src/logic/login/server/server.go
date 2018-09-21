package server

import (
	"context"
	"lib/comm/errors"
	"lib/comm/token"
	"log"

	comm "lib/comm/proto"
	proto "logic/login/proto"
)

type LoginHandler struct{}

func (handler *LoginHandler) Echo(ctx context.Context, req *proto.EchoRequest, rsp *proto.EchoResponse) error {
	rsp.Msg = req.Msg
	return nil
}

func (handler *LoginHandler) Login(ctx context.Context, req *comm.LoginRequest, rsp *comm.LoginResponse) error {
	log.Printf("req:%v", req)
	uid, err := token.CheckLoginToken(req.Token)
	if err != nil {
		log.Printf("CheckLoginToken err %v", err)
		return errors.LOGIC_LOGIN_TOKEN_ERROR
	}
	rsp.Uid = uid
	return nil
}
