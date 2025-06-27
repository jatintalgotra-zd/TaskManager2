package user

import (
	"database/sql"

	"TaskManager2/models"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) *store {
	return &store{db: db}
}

func (s *store) Create(u *models.User) (int64, error) {
	res, err := s.db.Exec("INSERT INTO users (name, email) VALUES ( ?, ?)", u.Name, u.Email)
	if err != nil {
		return 0, err
	}

	id, err2 := res.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

func (s *store) GetByID(id int64) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var u models.User

	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return &models.User{}, err
	}

	return &u, nil
}
