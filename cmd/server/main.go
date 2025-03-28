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

	qbuilder := sq.Insert("users").PlaceholderFormat(sq.Dollar).
		Columns("username", "email", "pass_hash", "user_role", "created_at", "updated_at").
		Values(in.Name, in.Email, in.Password, in.Role, time.Now(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := qbuilder.ToSql()
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

	qbuilder := sq.Select("id", "username", "email", "user_role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": in.GetId()}).
		Limit(1)

	query, args, err := qbuilder.ToSql()
	if err != nil {
		log.Fatal("Can't create query")
	}

	var id int64
	var name, email string
	var role generated.Role
	var created_at, updated_at time.Time

	err = serv.pool.QueryRow(context, query, args...).Scan(&id, &name, &email, &role, &created_at, &updated_at)
	if err != nil {
		return &generated.GetResponse{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: timestamppb.New(created_at),
			UpdatedAt: timestamppb.New(updated_at),
		}, err
	}

	//return some random data for testing
	res := &generated.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(created_at),
		UpdatedAt: timestamppb.New(updated_at),
	}

	return res, nil
}

func (serv *server) Update(context context.Context, in *generated.UpdateRequest) (*generated.UpdateResponse, error) {

	// Nullable values does not processing

	qbuilder := sq.Update("users").
		Set("username", in.Name.GetValue()).
		Set("email", in.Email.GetValue()).
		Where(sq.Eq{"id": in.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qbuilder.ToSql()
	if err != nil {
		log.Fatal("Can't create query")
	}

	_, err = serv.pool.Query(context, query, args...)
	if err != nil {
		return &generated.UpdateResponse{}, err
	}

	return &generated.UpdateResponse{}, nil
}

func (serv *server) Delete(context context.Context, in *generated.DeleteRequest) (*generated.DeleteResponse, error) {

	qbuilder := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": in.GetId()})

	query, args, err := qbuilder.ToSql()
	if err != nil {
		log.Fatal("Can't create query")
	}

	_, err = serv.pool.Query(context, query, args...)

	if err != nil {
		return &generated.DeleteResponse{}, err
	}

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
