package main

import (
	"log"
	"time"

	libMicro "lib/micro"
	"lib/pprof"
	loginProto "logic/login/proto"
	loginServer "logic/login/server"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
)

func main() {
	cmd.Init(
		cmd.Name(loginServer.ServiceName),
		cmd.Version(loginServer.ServiceVersion),
	)

	service := micro.NewService(
		micro.Name(loginServer.ServiceName),
		micro.Version(loginServer.ServiceVersion),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*10),
	)
	service.Init()
	service.Server().Init(
		server.WrapHandler(libMicro.ServerCommWrapper),
	)
	loginProto.RegisterLoginHandler(service.Server(), new(loginServer.LoginHandler))
	pprof.Init(loginServer.ServiceName, loginServer.ServiceVersion)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
