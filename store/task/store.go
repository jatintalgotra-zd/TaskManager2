package task

import (
	"errors"

	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

var errNotFound = errors.New("task not found")

type store struct {
}

func New() *store {
	return &store{}
}

func (store) Create(ctx *gofr.Context, t *models.Task) (int64, error) {
	db := ctx.SQL

	res, err := db.Exec("INSERT INTO tasks (description, status, user_id) VALUES ( ?, ?, ?)", t.Desc, t.Status, t.UserID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (store) GetAll(ctx *gofr.Context) ([]models.Task, error) {
	db := ctx.SQL

	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.UserID)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tasks, nil
}

func (store) GetByID(ctx *gofr.Context, id int64) (*models.Task, error) {
	db := ctx.SQL
	row := db.QueryRow("SELECT id, description, status, user_id FROM tasks WHERE id = ?", id)

	var t models.Task

	err := row.Scan(&t.ID, &t.Desc, &t.Status, &t.UserID)
	if err != nil {
		return &models.Task{}, err
	}

	return &t, nil
}

func (store) Update(ctx *gofr.Context, t *models.Task) error {
	db := ctx.SQL

	res, err := db.Exec("UPDATE tasks SET description = ?, status = ? WHERE id = ?", t.Desc, t.Status, t.ID)
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

func (store) Delete(ctx *gofr.Context, id int64) error {
	db := ctx.SQL

	res, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
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
