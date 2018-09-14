package main

import (
	"log"
	"time"

	"lib/pprof"
	proto "template/proto"
	"template/server"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
)

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
	proto.RegisterTemplateHandler(service.Server(), new(server.TemplateHandler))
	pprof.Init(server.ServiceName, server.ServiceVersion)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
