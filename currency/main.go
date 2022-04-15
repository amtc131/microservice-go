package main

import (
	"net"

	protos "github.com/amtc131/microservice-go/currency/protos/currency"
	"github.com/amtc131/microservice-go/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
	}

	gs.Serve(l)
}
