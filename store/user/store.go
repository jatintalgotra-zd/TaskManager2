package user

import (
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

type store struct {
}

func New() *store {
	return &store{}
}

func (store) Create(ctx *gofr.Context, u *models.User) (int64, error) {
	db := ctx.SQL

	res, err := db.Exec("INSERT INTO users (name, email) VALUES ( ?, ?)", u.Name, u.Email)
	if err != nil {
		return 0, err
	}

	id, err2 := res.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

func (store) GetByID(ctx *gofr.Context, id int64) (*models.User, error) {
	db := ctx.SQL
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var u models.User

	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return &models.User{}, err
	}

	return &u, nil
}
