package user

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"

	"TaskManager2/models"
	"TaskManager2/utils"
)

type lastInsertIDErrorResult struct{}

func (lastInsertIDErrorResult) LastInsertId() (int64, error) {
	return 0, utils.ErrTest
}

func (lastInsertIDErrorResult) RowsAffected() (int64, error) {
	return 1, nil
}

func TestStore_Create(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	userStore := New()
	query := "INSERT INTO users (name, email) VALUES ( ?, ?)"

	testcases := []struct {
		description   string
		input         *models.User
		mockExpect    func()
		wantID        int64
		expectedError bool
	}{
		{
			description: "success",
			input:       &models.User{Name: "test", Email: "test@example.com"},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("test", "test@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantID:        1,
			expectedError: false,
		},
		{
			description: "exec error",
			input:       &models.User{Name: "fail", Email: "fail@example.com"},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).WithArgs("fail", "fail@example.com").WillReturnError(utils.ErrTest)
			},
			wantID:        0,
			expectedError: true,
		},
		{
			description: "last inserted error",
			input:       &models.User{Name: "test", Email: "test@example.com"},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).WithArgs("test", "test@example.com").WillReturnResult(lastInsertIDErrorResult{})
			},
			wantID:        0,
			expectedError: true,
		},
	}

	for _, tc := range testcases {
		tc.mockExpect()

		id, err := userStore.Create(ctx, tc.input)
		if (err != nil) != tc.expectedError {
			t.Errorf("expected err: %v, got: %v", tc.expectedError, err)
		}

		if id != tc.wantID {
			t.Errorf("expected id: %v, got: %v", tc.wantID, id)
		}
	}
}

func TestStore_GetByID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	userStore := New()
	query := "SELECT id, name, email FROM users WHERE id = ?"

	testcases := []struct {
		description   string
		inputID       int64
		mockExpect    func()
		want          *models.User
		expectedError bool
	}{
		{
			description: "success",
			inputID:     1,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "test", "test@example.com")
				mock.SQL.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want:          &models.User{ID: 1, Name: "test", Email: "test"},
			expectedError: false,
		},
		{
			description: "scan error",
			inputID:     1,
			mockExpect: func() {
				rows2 := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")
				mock.SQL.ExpectQuery(query).WillReturnRows(rows2)
			},
			want:          &models.User{},
			expectedError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockExpect()

			user, err := userStore.GetByID(ctx, tc.inputID)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected err: %v, got: %v", tc.expectedError, err)
			}

			if user.ID != tc.want.ID {
				t.Errorf("expected id: %v, got: %v", tc.want.ID, user.ID)
			}
		})
	}
}
