package app

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	generated.UnimplementedAuthV1Server
	serviceProvider *serviceProvider
	server          *grpc.Server
	httpServer      *http.Server
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

	defer app.StopServices()
	return app.RunServices()

}

func (app *App) RunServices() error {

	starters := []func(*sync.WaitGroup){
		app.runGRPCServer,
		app.runHttpServer,
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(starters))

	for _, starter := range starters {
		go starter(&waitGroup)
	}

	waitGroup.Wait()
	return nil

}

func (app *App) StopServices() {
	app.serviceProvider.service.StopService(app.ctx)
	app.server.Stop()
}

func (app *App) initDeps(ctx context.Context) error {

	inits := []func(context.Context) error{
		app.initContext,
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
		app.initHttpServer,
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

func (app *App) initHttpServer(ctx context.Context) error {

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := generated.RegisterAuthV1HandlerFromEndpoint(ctx, mux, app.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	app.httpServer = &http.Server{
		Addr:    app.serviceProvider.HttpConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil

}

func (app *App) runGRPCServer(group *sync.WaitGroup) {

	defer group.Done()

	address := app.serviceProvider.GRPCConfig().Address()
	log.Printf("Server started on %v", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Print(err)
		return
	}

	err = app.server.Serve(listener)
	if err != nil {
		log.Print(err)
	}
}

func (app *App) runHttpServer(group *sync.WaitGroup) {

	defer group.Done()

	log.Printf("HTTP server is running on %s", app.serviceProvider.HttpConfig().Address())

	err := app.httpServer.ListenAndServe()

	if err != nil {
		log.Print(err)
	}

}
