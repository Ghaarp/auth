package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/Ghaarp/auth/internal/config"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"

	"github.com/Ghaarp/auth/internal/repository"
	repositoryInstance "github.com/Ghaarp/auth/internal/repository/auth"
	"github.com/Ghaarp/auth/internal/repository/auth/converter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	generated.UnimplementedAuthV1Server
	Repository repository.AuthRepository
	Converter  repository.RepoConverter
}

func (serv *server) Create(ctx context.Context, in *generated.CreateRequest) (*generated.CreateResponse, error) {

	privateData := &generated.PrivateUser{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     in.Role,
	}

	id, err := serv.Repository.Create(ctx, serv.Converter.ToRepoUserDataPrivate(privateData))

	return &generated.CreateResponse{
		Id: id,
	}, err
}

func (serv *server) Get(ctx context.Context, in *generated.GetRequest) (*generated.GetResponse, error) {

	userData, err := serv.Repository.Get(ctx, in.Id)
	if err != nil {
		return &generated.GetResponse{}, err
	}

	userDataProto := serv.Converter.ToProtoUserDataPublic(userData)
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

	err := serv.Repository.Update(ctx, serv.Converter.ToRepoUserDataPublic(userDataPublic))

	return &generated.UpdateResponse{}, err
}

func (serv *server) Delete(ctx context.Context, in *generated.DeleteRequest) (*generated.DeleteResponse, error) {

	err := serv.Repository.Delete(ctx, in.Id)
	return &generated.DeleteResponse{}, err
}

func main() {

	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Print("Unable to load .env")
	}

	authConfig, err := config.NewAuthConfig()
	if err != nil {
		log.Fatal("Unable to load auth config")
	}

	serv := &server{}
	addRepositoryLayer(serv, ctx)
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}

	defer serv.Repository.ClosePool(ctx)
	turnOnServer(serv, authConfig)
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

func addRepositoryLayer(serv *server, ctx context.Context) error {

	dbconfig, err := config.NewDBConfig()
	if err != nil {
		log.Fatal("Unable to load DB config")
	}

	serv.Repository, err = repositoryInstance.CreateRepository(ctx, dbconfig.DSN())
	if err != nil {
		return err
	}

	serv.Converter = &converter.AuthConverter{}
	return nil
}
