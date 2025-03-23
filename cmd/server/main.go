package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/Ghaarp/auth/internal/config"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4/pgxpool"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	generated.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

func (serv *server) Create(context context.Context, in *generated.CreateRequest) (*generated.CreateResponse, error) {

	builder := sq.Insert("users").PlaceholderFormat(sq.Dollar).
		Columns("username", "email", "pass_hash", "user_role", "created_at", "updated_at").
		Values(in.Name, in.Email, in.Password, in.Role, time.Now(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatal("Can't create query")
	}

	var userid int64
	err = serv.pool.QueryRow(context, query, args...).Scan(&userid)
	if err != nil {
		log.Fatal("Failed to create new user")
	}

	log.Printf("Created new user: %d", userid)

	return &generated.CreateResponse{
		Id: userid,
	}, nil
}

func (serv *server) Get(context context.Context, in *generated.GetRequest) (*generated.GetResponse, error) {
	log.Printf(color.GreenString("%v", in))

	//return some random data for testing
	res := &generated.GetResponse{
		Id:        gofakeit.Int64(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      generated.Role_R_USER,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	return res, nil
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

	dbconfig, err := config.NewDBConfig()
	if err != nil {
		log.Fatal("Unable to load DB config")
	}

	pool, err := pgxpool.Connect(ctx, dbconfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	turnOnServer(authConfig, pool)

}

func turnOnServer(conf config.AuthConfig, pool *pgxpool.Pool) {

	listener, err := net.Listen("tcp", conf.Address())
	if err != nil {
		log.Fatal(err)
	}

	serverObj := grpc.NewServer()
	reflection.Register(serverObj)
	serv := &server{}
	generated.RegisterAuthV1Server(serverObj, serv)
	serv.pool = pool
	log.Printf("Server started on %v", listener.Addr())

	if err := serverObj.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
