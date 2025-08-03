package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Ghaarp/auth/internal/repository/auth/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func CreateRepository(ctx context.Context, dsn string) (*repo, error) {
	repository := &repo{}
	err := repository.openPool(ctx, dsn)
	return repository, err
}

func (rep *repo) Create(ctx context.Context, data *model.UserDataPrivate) (int64, error) {
	qbuilder := sq.Insert("users").PlaceholderFormat(sq.Dollar).
		Columns("username", "email", "pass_hash", "user_role", "created_at", "updated_at").
		Values(data.Name, data.Email, data.Password, data.Role, time.Now(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := qbuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var userid int64
	err = rep.pool.QueryRow(ctx, query, args...).Scan(&userid)
	if err != nil {
		return 0, err
	}

	log.Printf("Created new user: %d", userid)

	return userid, nil
}

func (rep *repo) Get(ctx context.Context, id int64) (*model.UserDataPublic, error) {
	qbuilder := sq.Select("username", "email", "user_role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := qbuilder.ToSql()
	if err != nil {
		return &model.UserDataPublic{}, err
	}

	var name, email sql.NullString
	var role int64
	var created_at, updated_at time.Time

	err = rep.pool.QueryRow(ctx, query, args...).Scan(&name, &email, &role, &created_at, &updated_at)
	if err != nil {
		return &model.UserDataPublic{}, err
	}

	res := &model.UserDataPublic{
		Id:    id,
		Name:  name,
		Email: email,
		Role:  role,
	}

	return res, nil
}

func (rep *repo) Update(ctx context.Context, data *model.UserDataPublic) error {

	qbuilder := sq.Update("users").
		Set("username", data.Name).
		Set("email", data.Email).
		Where(sq.Eq{"id": data.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qbuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = rep.pool.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (rep *repo) Delete(ctx context.Context, id int64) error {

	qbuilder := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := qbuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = rep.pool.Query(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil

}

func (rep *repo) openPool(ctx context.Context, dsn string) error {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	rep.pool = pool
	return nil
}

func (rep *repo) ClosePool(ctx context.Context) {
	if rep.pool != nil {
		rep.pool.Close()
	}
}
