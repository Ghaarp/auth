package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ghaarp/auth/internal/repository/auth/model"
	sq "github.com/Masterminds/squirrel"
)

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
