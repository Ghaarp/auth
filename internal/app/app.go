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
