package task

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"TaskManager2/models"
	"TaskManager2/utils"
)

func TestStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock db: %s", err)
		return
	}

	defer db.Close()

	taskStore := New(db)
	query := "INSERT INTO tasks (description, status, user_id) VALUES ( ?, ?, ?)"

	// testcase 1 - success
	mock.ExpectExec(query).WithArgs("", false, 0).WillReturnResult(sqlmock.NewResult(1, 1))

	id, err2 := taskStore.Create(&models.Task{})
	if err2 != nil {
		t.Errorf("error creating task: %s", err2)
		return
	}

	if id != 1 {
		t.Errorf("want %d, got %d", 1, id)
		return
	}

	// testcase 2 - exec error
	mock.ExpectExec(query).WithArgs("fail", false, 0).WillReturnError(utils.ErrTest)

	_, err3 := taskStore.Create(&models.Task{Desc: "fail"})
	if err3 == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// testcase 3 - result err
	mock.ExpectExec(query).WithArgs("", false, 0).WillReturnResult(sqlmock.NewErrorResult(utils.ErrTest))

	_, err4 := taskStore.Create(&models.Task{Desc: "fail"})
	if err4 == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestStore_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock db: %s", err)
		return
	}

	defer db.Close()
	taskStore := New(db)

	query := "SELECT * FROM tasks"

	// testcase 1 - success
	rows := sqlmock.NewRows([]string{"id", "desc", "status", "user_id"}).AddRow(1, "test", false, "1")

	mock.ExpectQuery(query).WillReturnRows(rows)

	tasks, err2 := taskStore.GetAll()
	if err2 != nil {
		t.Errorf("error getting tasks: %s", err2)
		return
	}

	if len(tasks) != 1 {
		t.Errorf("want 1, got %d", len(tasks))
		return
	}

	if tasks[0].ID != 1 {
		t.Errorf("want 1, got %d", tasks[0].ID)
		return
	}

	// testcase 2 - query error
	mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

	_, err3 := taskStore.GetAll()
	if err3 == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// testcase 3 - scan error
	rows = sqlmock.NewRows([]string{"id", "desc", "status"}).AddRow(1, "test", false)

	mock.ExpectQuery(query).WillReturnRows(rows)

	_, err2 = taskStore.GetAll()
	if err2 == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestStore_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock db: %s", err)
		return
	}
	defer db.Close()

	taskStore := New(db)
	query := "SELECT id, description, status, user_id FROM tasks WHERE id = ?"
	rows := sqlmock.NewRows([]string{"id", "desc", "status", "user_id"}).AddRow(1, "test", false, "1")

	// testcase 1 - success

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	task, err2 := taskStore.GetByID(1)
	if err2 != nil {
		t.Errorf("error getting task: %s", err2)
		return
	}

	if task.ID != 1 {
		t.Errorf("want 1, got %d", task.ID)
		return
	}

	// testcase 2 - scan error
	rows2 := sqlmock.NewRows([]string{"id", "desc", "status"}).AddRow(1, "test", false)

	mock.ExpectQuery(query).WillReturnRows(rows2)

	_, err3 := taskStore.GetByID(1)
	if err3 == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock db: %s", err)
		return
	}
	defer db.Close()

	taskStore := New(db)
	query := "UPDATE tasks SET description = ?, status = ? WHERE id = ?"

	// testcase 1 - success
	mock.ExpectExec(query).WithArgs("test", true, int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	err2 := taskStore.Update(&models.Task{ID: 1, Desc: "test", Status: true})
	if err2 != nil {
		t.Errorf("error updating task: %s", err2)
		return
	}

	// testcase 2 - no rows affected
	mock.ExpectExec(query).WithArgs("test", false, int64(2)).WillReturnResult(sqlmock.NewResult(0, 0))

	err3 := taskStore.Update(&models.Task{ID: 2, Desc: "test", Status: false})
	if err3 == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// testcase 3 - exec error
	mock.ExpectExec(query).WithArgs("fail", false, int64(3)).WillReturnError(utils.ErrTest)

	err4 := taskStore.Update(&models.Task{ID: 3, Desc: "fail", Status: false})
	if err4 == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock db: %s", err)
		return
	}
	defer db.Close()

	taskStore := New(db)
	query := "DELETE FROM tasks WHERE id = ?"

	// testcase 1 - success
	mock.ExpectExec(query).WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	err2 := taskStore.Delete(1)
	if err2 != nil {
		t.Errorf("error deleting task: %s", err2)
		return
	}

	// testcase 2 - no rows affected
	mock.ExpectExec(query).WithArgs(int64(2)).WillReturnResult(sqlmock.NewResult(0, 0))

	err3 := taskStore.Delete(2)
	if err3 == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// testcase 3 - exec error
	mock.ExpectExec(query).WithArgs(int64(3)).WillReturnError(utils.ErrTest)

	err4 := taskStore.Delete(3)
	if err4 == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
