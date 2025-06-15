package app

import (
	"context"
	"flag"
	"log"
	"net"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	generated.UnimplementedAuthV1Server
	serviceProvider *serviceProvider
	server          *grpc.Server
	ctx             context.Context
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (app *App) Run() error {
	defer func() {
		app.serviceProvider.service.StopService(app.ctx)
		app.server.Stop()
	}()

	return app.runGRPCServer()
}

func (app *App) initDeps(ctx context.Context) error {

	inits := []func(context.Context) error{
		app.initContext,
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
	//defer app.serviceProvider.Service(ctx).StopService(ctx)
	//turnOnServer(serv, serv.serviceProvider.ServerConfig(configPath))
}

func (app *App) initContext(ctx context.Context) error {
	app.ctx = ctx
	return nil
}

func (app *App) initConfig(ctx context.Context) error {
	flag.Parse()
	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.server = grpc.NewServer()
	reflection.Register(app.server)
	generated.RegisterAuthV1Server(app.server, app.serviceProvider.AuthImplementation(ctx))
	return nil
}

func (app *App) runGRPCServer() error {

	listener, err := net.Listen("tcp", app.serviceProvider.ServerConfig(configPath).Address())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started on %v", listener.Addr())

	if err := app.server.Serve(listener); err != nil {
		log.Fatal(err)
	}
	return nil
}

/*
func turnOnServer(serv *server, conf config.AuthConfig) {

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
}*/
