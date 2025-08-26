package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"http-task-executor/task-executor/internal/task-executor/logger"
	"http-task-executor/task-executor/internal/task-executor/models"
	"strings"
)

// TaskRepository represents db repository to work with models.Task.
type TaskRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewRepository - creates new instance of TaskRepository.
func NewRepository(db *sqlx.DB, log logger.Logger) *TaskRepository {
	return &TaskRepository{db: db, log: log}
}

// UpdateResult updates task result fields.
func (r *TaskRepository) UpdateResult(ctx context.Context, task *models.Task) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return errors.Wrap(err, "TaskRepository.UpdateResult.BeginTx")
	}

	prepare, err := tx.PrepareContext(ctx, "UPDATE task SET status = $1, response_status_code = $2, response_length = $3 WHERE id = $4")
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return errors.Wrap(err1, "TaskRepository.UpdateResult.PrepareContext.Rollback")
		}
		return errors.Wrap(err, "TaskRepository.UpdateResult.PrepareContext")
	}
	res, err := prepare.ExecContext(ctx, task.Status, task.ResponseStatus, task.ResponseLength, task.ID)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return errors.Wrap(err1, "TaskRepository.UpdateResult.ExecContext.Rollback")
		}
		return errors.Wrap(err, "TaskRepository.UpdateResult.ExecContext")
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "TaskRepository.UpdateResult.RowsAffected")
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	outputHeaders := make([]models.Header, 0)
	for _, header := range task.Headers {
		if !header.Input {
			outputHeaders = append(outputHeaders, header)
		}
	}
	err = createHeaders(ctx, tx, task.ID, outputHeaders)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return errors.Wrap(err1, "TaskRepository.UpdateResult.createHeaders.Rollback")
		}
		return errors.Wrap(err, "TaskRepository.UpdateResult.createHeaders")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "TaskRepository.UpdateResult.Commit")
	}

	return nil
}

// UpdateStatus updates task status field.
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
