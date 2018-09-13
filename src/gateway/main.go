package main

import (
	"log"
	"time"

	proto "gateway/proto"
	"gateway/server"
	"lib/pprof"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	//_ "github.com/micro/go-plugins/registry/etcdv3"
)

var defaultGateway *server.Gateway

func main() {
	cmd.Init(
		cmd.Name(server.ServiceName),
		cmd.Version(server.ServiceVersion),
	)

	service := micro.NewService(
		micro.Name(server.ServiceName),
		micro.Version(server.ServiceVersion),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*10),
	)
	service.Init()
	proto.RegisterGatewayHandler(service.Server(), new(server.GatewayHandler))

	defaultGateway = server.NewGateway()
	if err := defaultGateway.Init(service); err != nil {
		log.Fatal(err)
	}
	defaultGateway.Run()

	pprof.Init(server.ServiceName, server.ServiceVersion)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
