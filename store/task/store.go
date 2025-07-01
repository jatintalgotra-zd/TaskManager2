package task

import (
	"database/sql"
	"errors"

	"TaskManager2/models"
)

var errNotFound = errors.New("task not found")

type store struct {
	db *sql.DB
}

func New(db *sql.DB) *store {
	return &store{db: db}
}

func (s *store) Create(t *models.Task) (int64, error) {
	res, err := s.db.Exec("INSERT INTO tasks (description, status, user_id) VALUES ( ?, ?, ?)", t.Desc, t.Status, t.UserID)
	if err != nil {
		return 0, err
	}

	id, err2 := res.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

func (s *store) GetAll() ([]models.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		err2 := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.UserID)
		if err2 != nil {
			return nil, err2
		}

		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tasks, nil
}

func (s *store) GetByID(id int64) (*models.Task, error) {
	row := s.db.QueryRow("SELECT id, description, status, user_id FROM tasks WHERE id = ?", id)

	var t models.Task

	err := row.Scan(&t.ID, &t.Desc, &t.Status, &t.UserID)
	if err != nil {
		return &models.Task{}, err
	}

	return &t, nil
}

func (s *store) Update(t *models.Task) error {
	res, err := s.db.Exec("UPDATE tasks SET description = ?, status = ? WHERE id = ?", t.Desc, t.Status, t.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errNotFound
	}

	return nil
}

func (s *store) Delete(id int64) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errNotFound
	}

	return nil
}
