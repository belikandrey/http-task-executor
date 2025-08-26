package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"task-service/internal/logger"
	"task-service/internal/models"
)

type TaskRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewRepository(db *sqlx.DB, log logger.Logger) *TaskRepository {
	return &TaskRepository{db: db, log: log}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, errors.Wrap(err, "TaskRepository.Create.BeginTx")
	}

	if task.Headers == nil {
		task.Headers = make([]models.Header, 0)
	}

	prepare, err := tx.PrepareContext(ctx, "INSERT INTO task (method, url, status, response_status_code, response_length) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return nil, errors.Wrap(err1, "TaskRepository.Create.PrepareContext.Rollback")
		}
		return nil, errors.Wrap(err, "TaskRepository.Create.PrepareContext")
	}
	var id int64
	rowContext := prepare.QueryRowContext(ctx, task.Method, task.Url, task.Status, task.ResponseStatus, task.ResponseLength)
	err = rowContext.Scan(&id)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return nil, errors.Wrap(err1, "TaskRepository.Create.QueryRowContext.Rollback")
		}
		return nil, errors.Wrap(err, "TaskRepository.Create.QueryRowContext")
	}

	err = createHeaders(ctx, tx, id, task.Headers)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return nil, errors.Wrap(err1, "TaskRepository.Create.createHeaders.Rollback")
		}
		return nil, errors.Wrap(err, "TaskRepository.Create.createHeaders")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "TaskRepository.Create.Commit")
	}
	task.Id = id

	return task, nil
}

func (r *TaskRepository) GetByIdWithOutputHeaders(ctx context.Context, id int64) (*models.Task, error) {
	prepareContext, err := r.db.PrepareContext(ctx, `SELECT t.id,
       								t.url as url,
       								t.method as method,
									t.status as status,
									t.response_status_code as response_status,
									t.response_length as response_length,
									COALESCE(h.name, '') as header_name,
									COALESCE(h.value, '') as header_value
									FROM task t
									LEFT JOIN headers h ON h.task_id = t.id AND h.input=false
									WHERE t.id = $1`)
	if err != nil {
		return nil, errors.Wrap(err, "TaskRepository.GetByIdWithResponseHeaders.PrepareContext")
	}
	rows, err := prepareContext.QueryContext(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "TaskRepository.GetByIdWithResponseHeaders.QueryContext")
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.Errorf("TaskRepository.GetByIdWithResponseHeaders.rows.Close(): %v", err)
		}
	}(rows)

	var task *models.Task
	tempTask := &models.Task{}
	for rows.Next() {
		header := models.Header{Input: false}
		if task == nil {
			task = &models.Task{}
			task.Headers = make([]models.Header, 0)
			err = rows.Scan(&task.Id, &task.Url, &task.Method, &task.Status, &task.ResponseStatus, &task.ResponseLength, &header.Name, &header.Value)
		} else {
			err = rows.Scan(&tempTask.Id, &tempTask.Url, &tempTask.Method, &tempTask.Status, &tempTask.ResponseStatus, &tempTask.ResponseLength, &header.Name, &header.Value)
		}
		if err != nil {
			return nil, err
		}
		if header.Name != "" && header.Value != "" {
			task.Headers = append(task.Headers, header)
		}
	}
	if task == nil {
		return nil, sql.ErrNoRows
	}

	return task, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id int64, newStatus string) error {
	prepareContext, err := r.db.PrepareContext(ctx, "UPDATE task SET status=$1 WHERE id=$2")
	if err != nil {
		return errors.Wrap(err, "TaskRepository.UpdateStatus.PrepareContext")
	}

	result, err := prepareContext.ExecContext(ctx, newStatus, id)
	if err != nil {
		return errors.Wrap(err, "TaskRepository.UpdateStatus.ExecContext")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "TaskRepository.UpdateStatus.RowsAffected")
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	prepareContext, err := r.db.PrepareContext(ctx, "DELETE FROM task WHERE id=$1")
	if err != nil {
		return errors.Wrap(err, "TaskRepository.Delete.PrepareContext")
	}

	result, err := prepareContext.ExecContext(ctx, id)
	if err != nil {
		return errors.Wrap(err, "TaskRepository.Delete.ExecContext")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "TaskRepository.Delete.RowsAffected")
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func createHeaders(ctx context.Context, tx *sql.Tx, taskId int64, headers []models.Header) error {
	if len(headers) == 0 {
		return nil
	}
	sb := new(strings.Builder)
	sb.WriteString("INSERT INTO headers(name, value, input, task_id) VALUES ")
	params := make([]interface{}, 0, len(headers)*2)
	counter := 1
	for _, v := range headers {
		separator := ","
		params = append(params, v.Name, v.Value, v.Input)
		_, err := fmt.Fprintf(sb, "($%d, $%d, $%d, %d) %s", counter, counter+1, counter+2, taskId, separator)
		if err != nil {
			return err
		}
		counter += 3
	}
	s := sb.String()
	s = s[:len(s)-1]
	prepare, err := tx.PrepareContext(ctx, s)
	if err != nil {
		return err
	}
	_, err = prepare.ExecContext(ctx, params...)
	if err != nil {
		return err
	}
	return nil
}
