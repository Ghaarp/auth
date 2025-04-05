package model

import (
	"database/sql"
)

type UserDataPrivate struct {
	Name     string
	Email    string
	Password string
	Role     int64
}

type UserDataPublic struct {
	Id    int64
	Name  sql.NullString
	Email sql.NullString
	Role  int64
}
