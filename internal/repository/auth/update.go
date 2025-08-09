package auth

import (
	"context"

	"github.com/Ghaarp/auth/internal/repository/auth/model"
	sq "github.com/Masterminds/squirrel"
)

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
