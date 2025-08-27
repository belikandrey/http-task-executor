package repository

import (
	"context"
	dbSql "database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"http-task-executor/task-service/internal/task-service/models"
	"testing"
)

func TestTasksRepo_CreateWithoutHeaders(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDb := sqlx.NewDb(db, "sqlmock")

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	tasksRepo := NewRepository(sqlxDb, sugar)

	t.Run("Create", func(t *testing.T) {
		task := &models.Task{
			Method: "GET",
			URL:    "https://www.google.com",
			Status: models.StatusNew,
		}

		sql := "INSERT INTO task (method, url, status, response_status_code, response_length) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(task.Method, task.URL, task.Status, task.ResponseStatus, task.ResponseLength).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		created, err := tasksRepo.Create(context.Background(), task)

		require.NoError(t, err)
		require.NotEmpty(t, created)
		assert.Equal(t, int64(1), created.ID)
		assert.Equal(t, task.Method, created.Method)
		assert.Equal(t, task.URL, created.URL)
		assert.Equal(t, task.Status, created.Status)
	})
}

func TestTasksRepo_CreateWithHeaders(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDb := sqlx.NewDb(db, "sqlmock")

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	tasksRepo := NewRepository(sqlxDb, sugar)
	headers := make([]models.Header, 0)

	header := models.Header{Name: "TEST_NAME", Value: "TEST_VALUE", Input: true}
	headers = append(headers, header)

	twoHeaders := make([]models.Header, 0)
	twoHeaders = append(twoHeaders, header)
	secondHeader := models.Header{Name: "TEST_NAME2", Value: "TEST_VALUE2", Input: true}
	twoHeaders = append(twoHeaders, secondHeader)

	t.Run("Create with one header", func(t *testing.T) {
		task := &models.Task{
			Method:  "GET",
			URL:     "https://www.google.com",
			Status:  models.StatusNew,
			Headers: headers,
		}

		sql := "INSERT INTO task (method, url, status, response_status_code, response_length) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(task.Method, task.URL, task.Status, task.ResponseStatus, task.ResponseLength).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		created, err := tasksRepo.Create(context.Background(), task)

		require.NoError(t, err)
		require.NotEmpty(t, created)
		assert.Equal(t, int64(1), created.ID)
		assert.Equal(t, task.Method, created.Method)
		assert.Equal(t, task.URL, created.URL)
		assert.Equal(t, task.Status, created.Status)
	})

	t.Run("Create with two headers", func(t *testing.T) {
		task := &models.Task{
			Method:  "GET",
			URL:     "https://www.google.com",
			Status:  models.StatusNew,
			Headers: twoHeaders,
		}

		sql := "INSERT INTO task (method, url, status, response_status_code, response_length) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1) ,($4, $5, $6, 1) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(task.Method, task.URL, task.Status, task.ResponseStatus, task.ResponseLength).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input, secondHeader.Name, secondHeader.Value, secondHeader.Input).WillReturnResult(sqlmock.NewResult(1, 2))
		mock.ExpectCommit()

		created, err := tasksRepo.Create(context.Background(), task)

		require.NoError(t, err)
		require.NotEmpty(t, created)
		assert.Equal(t, int64(1), created.ID)
		assert.Equal(t, task.Method, created.Method)
		assert.Equal(t, task.URL, created.URL)
		assert.Equal(t, task.Status, created.Status)
	})

	t.Run("Expect transaction rollback if cannot create headers", func(t *testing.T) {
		task := &models.Task{
			Method:  "GET",
			URL:     "https://www.google.com",
			Status:  models.StatusNew,
			Headers: twoHeaders,
		}

		sql := "INSERT INTO task (method, url, status, response_status_code, response_length) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		headersSql := "INSERT INTO headers(name, value, input, task_id) VALUES ($1, $2, $3, 1) ,($4, $5, $6, 1) "
		mock.ExpectBegin()
		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(task.Method, task.URL, task.Status, task.ResponseStatus, task.ResponseLength).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectPrepare(headersSql)
		mock.ExpectExec(headersSql).WithArgs(header.Name, header.Value, header.Input, secondHeader.Name, secondHeader.Value, secondHeader.Input).WillReturnError(errors.New("error"))
		mock.ExpectRollback()

		created, err := tasksRepo.Create(context.Background(), task)

		require.Error(t, err)
		require.Empty(t, created)
	})
}

func TestTasksRepo_GetByIdWithOutputHeaders(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDb := sqlx.NewDb(db, "sqlmock")

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	tasksRepo := NewRepository(sqlxDb, sugar)

	sql := `SELECT t.id,
       								t.url as url,
       								t.method as method,
									t.status as status,
									t.response_status_code as response_status,
									t.response_length as response_length,
									COALESCE(h.name, '') as header_name,
									COALESCE(h.value, '') as header_value
									FROM task t
									LEFT JOIN headers h ON h.task_id = t.id AND h.input=false
									WHERE t.id = $1`

	t.Run("GetById with one header", func(t *testing.T) {
		id := int64(1)
		url := "https://www.google.com"
		method := "GET"
		status := models.StatusNew
		responseStatusCode := int64(200)
		responseLength := int64(10)
		headerName := "TEST_NAME"
		headerValue := "TEST_VALUE"

		rows := sqlmock.NewRows([]string{"id", "url", "method", "status", "response_status_code", "response_length", "header_name", "header_value"}).AddRow(id, url, method, status, responseStatusCode, responseLength, headerName, headerValue)

		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(id).WillReturnRows(rows)

		task, err := tasksRepo.GetByIDWithOutputHeaders(context.Background(), id)

		require.NoError(t, err)
		require.NotEmpty(t, task)
		assert.Equal(t, id, task.ID)
		assert.Equal(t, method, task.Method)
		assert.Equal(t, url, task.URL)
		assert.Equal(t, status, task.Status)
		assert.NotEmpty(t, task.ResponseStatus)
		assert.NotEmpty(t, task.ResponseLength)
		assert.Equal(t, responseStatusCode, *task.ResponseStatus)
		assert.Equal(t, responseLength, *task.ResponseLength)
		assert.NotEmpty(t, task.Headers)
		assert.Equal(t, headerName, task.Headers[0].Name)
		assert.Equal(t, headerValue, task.Headers[0].Value)
	})

	t.Run("GetById without headers", func(t *testing.T) {
		id := int64(1)
		url := "https://www.google.com"
		method := "GET"
		status := models.StatusNew
		responseStatusCode := int64(200)
		responseLength := int64(10)

		rows := sqlmock.NewRows([]string{"id", "url", "method", "status", "response_status_code", "response_length", "header_name", "header_value"}).AddRow(id, url, method, status, responseStatusCode, responseLength, "", "")

		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(id).WillReturnRows(rows)

		task, err := tasksRepo.GetByIDWithOutputHeaders(context.Background(), id)

		require.NoError(t, err)
		require.NotEmpty(t, task)
		assert.Equal(t, id, task.ID)
		assert.Equal(t, method, task.Method)
		assert.Equal(t, url, task.URL)
		assert.Equal(t, status, task.Status)
		assert.NotEmpty(t, task.ResponseStatus)
		assert.NotEmpty(t, task.ResponseLength)
		assert.Equal(t, responseStatusCode, *task.ResponseStatus)
		assert.Equal(t, responseLength, *task.ResponseLength)
		assert.Empty(t, task.Headers)
	})

	t.Run("GetById with 2 headers", func(t *testing.T) {
		id := int64(1)
		url := "https://www.google.com"
		method := "GET"
		status := models.StatusNew
		responseStatusCode := int64(200)
		responseLength := int64(10)
		headerName := "TEST_NAME"
		headerValue := "TEST_VALUE"
		headerName2 := "TEST_NAME2"
		headerValue2 := "TEST_VALUE2"

		rows := sqlmock.NewRows([]string{"id", "url", "method", "status", "response_status_code", "response_length", "header_name", "header_value"}).
			AddRow(id, url, method, status, responseStatusCode, responseLength, headerName, headerValue).
			AddRow(id, url, method, status, responseStatusCode, responseLength, headerName2, headerValue2)

		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(id).WillReturnRows(rows)

		task, err := tasksRepo.GetByIDWithOutputHeaders(context.Background(), id)

		require.NoError(t, err)
		require.NotEmpty(t, task)
		assert.Equal(t, id, task.ID)
		assert.Equal(t, method, task.Method)
		assert.Equal(t, url, task.URL)
		assert.Equal(t, status, task.Status)
		assert.NotEmpty(t, task.ResponseStatus)
		assert.NotEmpty(t, task.ResponseLength)
		assert.Equal(t, responseStatusCode, *task.ResponseStatus)
		assert.Equal(t, responseLength, *task.ResponseLength)
		assert.NotEmpty(t, task.Headers)
		assert.Len(t, task.Headers, 2)
	})

	t.Run("GetById with empty result", func(t *testing.T) {
		id := int64(1515)

		rows := sqlmock.NewRows([]string{"id", "url", "method", "status", "response_status_code", "response_length", "header_name", "header_value"})

		mock.ExpectPrepare(sql)
		mock.ExpectQuery(sql).WithArgs(id).WillReturnRows(rows)

		task, err := tasksRepo.GetByIDWithOutputHeaders(context.Background(), id)

		require.Error(t, err)
		require.Nil(t, task)
		require.ErrorIs(t, err, dbSql.ErrNoRows)
	})
}

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

func TestTasksRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDb := sqlx.NewDb(db, "sqlmock")

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	sql := "DELETE FROM task WHERE id=$1"

	tasksRepo := NewRepository(sqlxDb, sugar)

	t.Run("Delete successfully", func(t *testing.T) {
		id := int64(1515)

		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

		err := tasksRepo.Delete(context.Background(), id)

		require.NoError(t, err)
	})

	t.Run("Delete not rows affected", func(t *testing.T) {
		id := int64(1515)

		mock.ExpectPrepare(sql)
		mock.ExpectExec(sql).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 0))

		err := tasksRepo.Delete(context.Background(), id)

		require.Error(t, err)
		require.ErrorIs(t, err, dbSql.ErrNoRows)
	})
}
