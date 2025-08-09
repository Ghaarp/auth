package auth

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

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
