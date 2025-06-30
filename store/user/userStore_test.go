package user

import (
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

	userStore := New(db)
	query := "INSERT INTO users (name, email) VALUES ( ?, ?)"

	// testcase 1 - success
	mock.ExpectExec(query).WithArgs("test", "test@example.com").WillReturnResult(sqlmock.NewResult(1, 1))

	id, err2 := userStore.Create(&models.User{Name: "test", Email: "test@example.com"})
	if err2 != nil {
		t.Errorf("error creating user: %s", err2)
		return
	}

	if id != 1 {
		t.Errorf("want 1, got %d", id)
		return
	}

	// testcase 2 - exec error
	mock.ExpectExec(query).WithArgs("fail", "fail@example.com").WillReturnError(utils.ErrTest)

	_, err3 := userStore.Create(&models.User{Name: "fail", Email: "fail@example.com"})
	if err3 == nil {
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

	userStore := New(db)
	query := "SELECT id, name, email FROM users WHERE id = ?"

	// testcase 1 - success
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "test", "test@example.com")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	user, err2 := userStore.GetByID(1)
	if err2 != nil {
		t.Errorf("error getting user: %s", err2)
		return
	}

	if user.ID != 1 {
		t.Errorf("want 1, got %d", user.ID)
		return
	}

	// testcase 2 - scan error
	rows2 := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")

	mock.ExpectQuery(query).WillReturnRows(rows2)

	_, err3 := userStore.GetByID(1)
	if err3 == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
