package task

import (
	"database/sql"
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

type rowsAffectedErrorResult struct{}

func (rowsAffectedErrorResult) RowsAffected() (int64, error) {
	return 0, utils.ErrTest
}

func (rowsAffectedErrorResult) LastInsertId() (int64, error) {
	return 1, nil
}

func TestStore_Create(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	taskStore := New()
	query := "INSERT INTO tasks (description, status, user_id) VALUES ( ?, ?, ?)"

	tests := []struct {
		description   string
		input         *models.Task
		mockExpect    func()
		wantID        int64
		expectedError bool
	}{
		{
			description: "success",
			input:       &models.Task{},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("", false, 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantID:        1,
			expectedError: false,
		},
		{
			description: "exec error",
			input:       &models.Task{Desc: ""},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("", false, 0).
					WillReturnError(utils.ErrTest)
			},
			expectedError: true,
		},
		{
			description: "lastInsertID error",
			input:       &models.Task{},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("", false, 0).
					WillReturnResult(lastInsertIDErrorResult{})
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockExpect()

			id, err := taskStore.Create(ctx, tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected err: %v, got: %v", tc.expectedError, err)
			}

			if id != tc.wantID {
				t.Errorf("expected id: %d, got: %d", tc.wantID, id)
			}
		})
	}
}

func TestStore_GetAll(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	taskStore := New()
	query := "SELECT * FROM tasks"

	tests := []struct {
		description   string
		mockExpect    func()
		wantLen       int
		expectedError bool
	}{
		{
			description: "success",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "desc", "status", "user_id"}).
					AddRow(1, "test", false, "1")
				mock.SQL.ExpectQuery(query).WillReturnRows(rows)
			},
			wantLen:       1,
			expectedError: false,
		},
		{
			description: "query error",
			mockExpect: func() {
				mock.SQL.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
			},
			wantLen:       0,
			expectedError: true,
		},
		{
			description: "scan error",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "desc", "status"}).
					AddRow(1, "test", false)
				mock.SQL.ExpectQuery(query).WillReturnRows(rows)
			},
			wantLen:       0,
			expectedError: true,
		},
		{
			description: "row error",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "desc", "status", "user_id"}).
					AddRow(1, "test", false, "1").
					RowError(0, utils.ErrTest)
				mock.SQL.ExpectQuery(query).WillReturnRows(rows)
			},
			wantLen:       0,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockExpect()

			tasks, err := taskStore.GetAll(ctx)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error = %v, got = %v", tc.expectedError, err)
			}

			if len(tasks) != tc.wantLen {
				t.Errorf("expected task count = %d, got = %d", tc.wantLen, len(tasks))
			}
		})
	}
}

func TestStore_GetByID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	taskStore := New()
	query := "SELECT id, description, status, user_id FROM tasks WHERE id = ?"

	tests := []struct {
		description   string
		inputID       int64
		mockExpect    func()
		want          *models.Task
		expectedError bool
	}{
		{
			description: "success",
			inputID:     1,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "desc", "status", "user_id"}).
					AddRow(1, "test", false, "1")
				mock.SQL.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want:          &models.Task{ID: 1},
			expectedError: false,
		},
		{
			description: "scan error - missing user_id",
			inputID:     1,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "desc", "status"}).
					AddRow(1, "test", false)
				mock.SQL.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want:          nil,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockExpect()

			got, err := taskStore.GetByID(ctx, tc.inputID)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error = %v, got error = %v", tc.expectedError, err)
			}

			if !tc.expectedError && got.ID != tc.want.ID {
				t.Errorf("expected ID = %d, got = %d", tc.want.ID, got.ID)
			}
		})
	}
}

func TestStore_Update(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	taskStore := New()
	query := "UPDATE tasks SET description = ?, status = ? WHERE id = ?"

	tests := []struct {
		description   string
		input         *models.Task
		mockExpect    func()
		expectedError bool
	}{
		{
			description: "success",
			input:       &models.Task{ID: 1, Desc: "test", Status: true},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("test", true, int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			description: "no rows affected",
			input:       &models.Task{ID: 2, Desc: "test", Status: false},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("test", false, int64(2)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: true,
		},
		{
			description: "exec error",
			input:       &models.Task{ID: 3, Desc: "fail", Status: false},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("fail", false, int64(3)).
					WillReturnError(utils.ErrTest)
			},
			expectedError: true,
		},
		{
			description: "rowsAffected error",
			input:       &models.Task{ID: 1, Desc: "test", Status: true},
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs("test", true, int64(1)).
					WillReturnResult(rowsAffectedErrorResult{})
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockExpect()

			err := taskStore.Update(ctx, tc.input)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestStore_Delete(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	taskStore := New()
	query := "DELETE FROM tasks WHERE id = ?"

	tests := []struct {
		description   string
		inputID       int64
		mockExpect    func()
		expectedError bool
	}{
		{
			description: "success",
			inputID:     1,
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			description: "no rows affected",
			inputID:     2,
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs(int64(2)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: true,
		},
		{
			description: "exec error",
			inputID:     3,
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs(int64(3)).
					WillReturnError(utils.ErrTest)
			},
			expectedError: true,
		},
		{
			description: "rowsAffected error",
			inputID:     1,
			mockExpect: func() {
				mock.SQL.ExpectExec(query).
					WithArgs(int64(1)).
					WillReturnResult(rowsAffectedErrorResult{})
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockExpect()

			err := taskStore.Delete(ctx, tc.inputID)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
