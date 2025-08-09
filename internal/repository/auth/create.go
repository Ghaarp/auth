package auth

import (
	"context"
	"log"
	"time"

	"github.com/Ghaarp/auth/internal/repository/auth/model"
	sq "github.com/Masterminds/squirrel"
)

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
