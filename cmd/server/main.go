package main

import (
	"context"
	"fmt"
	"log"
	"net"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"
	"github.com/fatih/color"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	address = "localhost:"
	port    = 50051
)

type server struct {
	generated.UnimplementedAuthV1Server
}

func (serv *server) Create(context context.Context, in *generated.CreateRequest) (*generated.CreateResponse, error) {
	log.Printf(color.GreenString("%v", in))

	return &generated.CreateResponse{}, nil
}

func (serv *server) Get(context context.Context, in *generated.GetRequest) (*generated.GetResponse, error) {
	log.Printf(color.GreenString("%v", in))

	return &generated.GetResponse{}, nil
}

func (serv *server) Update(context context.Context, in *generated.UpdateRequest) (*generated.UpdateResponse, error) {
	log.Printf(color.GreenString("%v", in))

	return &generated.UpdateResponse{}, nil
}

func (serv *server) Delete(context context.Context, in *generated.DeleteRequest) (*generated.DeleteResponse, error) {
	log.Printf(color.GreenString("%v", in))

	return &generated.DeleteResponse{}, nil
}

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatal(err)
	}

	serverObj := grpc.NewServer()
	reflection.Register(serverObj)
	a := &server{}
	generated.RegisterAuthV1Server(serverObj, a)

	log.Printf("Server started on %v", listener.Addr())

	if err := serverObj.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
