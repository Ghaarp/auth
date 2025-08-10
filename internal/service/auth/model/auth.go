package model

type UserDataPrivate struct {
	Name     string
	Email    string
	Password string
	Role     int64
}

type UserDataPublic struct {
	Id    int64
	Name  string
	Email string
	Role  int64
}
