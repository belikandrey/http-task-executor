package repository

import (
	"context"
	dbSql "database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"http-task-executor/task-executor/internal/task-executor/models"
	"testing"
)

func TestTasksRepo_UpdateStatus(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDb := sqlx.NewDb(db, "sqlmock")

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	sql := "UPDATE task SET status=$1 WHERE id=$2"

	tasksRepo := NewRepository(sqlxDb, sugar)

	t.Run("UpdateStatus successfully", func(t *testing.T) {
		id := int64(1515)
		newStatus := models.StatusInProcess

		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(newStatus, id).WillReturnResult(sqlmock.NewResult(1, 1))

		err := tasksRepo.UpdateStatus(context.Background(), id, newStatus)

		require.NoError(t, err)
	})

	t.Run("UpdateStatus not rows affected", func(t *testing.T) {
		id := int64(1515)
		newStatus := models.StatusInProcess

		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(newStatus, id).WillReturnResult(sqlmock.NewResult(1, 0))

		err := tasksRepo.UpdateStatus(context.Background(), id, newStatus)

		require.Error(t, err)
		require.ErrorIs(t, err, dbSql.ErrNoRows)
	})
}
func TestTasksRepo_UpdateResultWithoutHeaders(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDb := sqlx.NewDb(db, "sqlmock")

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	tasksRepo := NewRepository(sqlxDb, sugar)

	sql := "UPDATE task SET status = $1, response_status_code = $2, response_length = $3 WHERE id = $4"

	t.Run("Update result without headers", func(t *testing.T) {
		status := int64(200)
		responseLength := int64(10)
		task := &models.Task{
			ID:             int64(1515),
			Status:         models.StatusDone,
			ResponseStatus: &status,
			ResponseLength: &responseLength,
		}

		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(task.Status, task.ResponseStatus, task.ResponseLength, task.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := tasksRepo.UpdateResult(context.Background(), task)

		require.NoError(t, err)
	})

	t.Run("Update result with one header", func(t *testing.T) {
		status := int64(200)
		responseLength := int64(10)
		header := models.Header{Name: "TEST_NAME", Value: "TEST_VALUE", Input: false}
		headers := []models.Header{header}
		task := &models.Task{
			ID:             int64(1515),
			Status:         models.StatusDone,
			ResponseStatus: &status,
			ResponseLength: &responseLength,
			Headers:        headers,
		}
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1515) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(task.Status, task.ResponseStatus, task.ResponseLength, task.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := tasksRepo.UpdateResult(context.Background(), task)

		require.NoError(t, err)
	})

	t.Run("Update result with 2 headers", func(t *testing.T) {
		status := int64(200)
		responseLength := int64(10)
		header := models.Header{Name: "TEST_NAME", Value: "TEST_VALUE", Input: false}
		secondHeader := models.Header{Name: "TEST_NAME2", Value: "TEST_VALUE2", Input: false}
		headers := []models.Header{header, secondHeader}
		task := &models.Task{
			ID:             int64(1515),
			Status:         models.StatusDone,
			ResponseStatus: &status,
			ResponseLength: &responseLength,
			Headers:        headers,
		}
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1515) ,($4, $5, $6, 1515) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(task.Status, task.ResponseStatus, task.ResponseLength, task.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input, secondHeader.Name, secondHeader.Value, secondHeader.Input).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := tasksRepo.UpdateResult(context.Background(), task)

		require.NoError(t, err)
	})

	t.Run("Rollback if error occurred", func(t *testing.T) {
		status := int64(200)
		responseLength := int64(10)
		header := models.Header{Name: "TEST_NAME", Value: "TEST_VALUE", Input: false}
		headers := []models.Header{header}
		task := &models.Task{
			ID:             int64(1515),
			Status:         models.StatusDone,
			ResponseStatus: &status,
			ResponseLength: &responseLength,
			Headers:        headers,
		}
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1515) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(task.Status, task.ResponseStatus, task.ResponseLength, task.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input).WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err := tasksRepo.UpdateResult(context.Background(), task)

		require.Error(t, err)
	})

	t.Run("Error if no rows affected", func(t *testing.T) {
		status := int64(200)
		responseLength := int64(10)
		header := models.Header{Name: "TEST_NAME", Value: "TEST_VALUE", Input: false}
		headers := []models.Header{header}
		task := &models.Task{
			ID:             int64(1515),
			Status:         models.StatusDone,
			ResponseStatus: &status,
			ResponseLength: &responseLength,
			Headers:        headers,
		}
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1515) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(task.Status, task.ResponseStatus, task.ResponseLength, task.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input).WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err := tasksRepo.UpdateResult(context.Background(), task)

		require.Error(t, err)
		require.ErrorIs(t, err, dbSql.ErrNoRows)
	})
}
