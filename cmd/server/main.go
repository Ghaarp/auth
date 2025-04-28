package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/Ghaarp/auth/internal/config"
	serviceDef "github.com/Ghaarp/auth/internal/service"
	serviceConverter "github.com/Ghaarp/auth/internal/service/auth/converter"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	generated.UnimplementedAuthV1Server
	serviceProvider  *serviceProvider
	serviceConverter serviceDef.ServiceConverter
}

func (serv *server) Create(ctx context.Context, in *generated.CreateRequest) (*generated.CreateResponse, error) {

	privateData := &generated.PrivateUser{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     in.Role,
	}

	id, err := serv.serviceProvider.Service(ctx).Create(ctx, serv.serviceConverter.ToServiceUserDataPrivate(privateData))

	return &generated.CreateResponse{
		Id: id,
	}, err
}

func (serv *server) Get(ctx context.Context, in *generated.GetRequest) (*generated.GetResponse, error) {

	userData, err := serv.serviceProvider.Service(ctx).Get(ctx, in.Id)
	if err != nil {
		return &generated.GetResponse{}, err
	}

	userDataProto := serv.serviceConverter.ToProtoUserDataPublic(userData)
	return &generated.GetResponse{
		Id:    userDataProto.Id,
		Name:  userDataProto.Name,
		Email: userDataProto.Email,
		Role:  userDataProto.Role,
	}, nil
}

func (serv *server) Update(ctx context.Context, in *generated.UpdateRequest) (*generated.UpdateResponse, error) {

	// Nullable values does not processing yet
	userDataPublic := &generated.PublicUser{
		Id:    in.Id,
		Name:  in.Name.GetValue(),
		Email: in.Email.GetValue(),
	}

	err := serv.serviceProvider.Service(ctx).Update(ctx, serv.serviceConverter.ToServiceUserDataPublic(userDataPublic))

	return &generated.UpdateResponse{}, err
}

func (serv *server) Delete(ctx context.Context, in *generated.DeleteRequest) (*generated.DeleteResponse, error) {

	err := serv.serviceProvider.Service(ctx).Delete(ctx, in.Id)
	return &generated.DeleteResponse{}, err
}

func main() {

	flag.Parse()
	ctx := context.Background()

	serv := &server{}
	serv.serviceProvider = newServiceProvider()
	serv.serviceConverter = serviceConverter.CreateConverter()

	defer serv.serviceProvider.Service(ctx).StopService(ctx)
	turnOnServer(serv, serv.serviceProvider.ServerConfig(configPath))
}

func turnOnServer(serv *server, conf config.AuthConfig) {

	listener, err := net.Listen("tcp", conf.Address())
	if err != nil {
		log.Fatal(err)
	}

	serverObj := grpc.NewServer()
	reflection.Register(serverObj)
	generated.RegisterAuthV1Server(serverObj, serv)
	log.Printf("Server started on %v", listener.Addr())

	if err := serverObj.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
