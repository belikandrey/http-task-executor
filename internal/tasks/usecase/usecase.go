package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"http-task-executor/internal/logger"
	"http-task-executor/internal/models"
	"http-task-executor/internal/tasks"
	"http-task-executor/pkg/errors/general/validation"
	httpErrors "http-task-executor/pkg/errors/http"
	"http-task-executor/pkg/utils"
)

type TaskUseCase struct {
	log  logger.Logger
	repo tasks.Repository
	exec tasks.Executor
}

func NewTaskUseCase(log logger.Logger, repo tasks.Repository, exec tasks.Executor) *TaskUseCase {
	return &TaskUseCase{log: log, repo: repo, exec: exec}
}

func (t *TaskUseCase) Create(ctx context.Context, task *models.Task) (*models.Task, error) {

	validationErrors := validateTask(ctx, task)
	if validationErrors != nil && len(validationErrors) > 0 {
		return nil, httpErrors.NewValidationError(validationErrors)
	}

	create, err := t.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	go t.exec.ExecuteTask(*create)

	return create, nil
}

func (t *TaskUseCase) GetByIdWithOutputHeaders(ctx context.Context, id int64) (*models.Task, error) {
	if id <= 0 {
		return nil, httpErrors.NewBadRequestError(errors.New("invalid id"))
	}

	task, err := t.repo.GetByIdWithOutputHeaders(ctx, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func validateTask(ctx context.Context, task *models.Task) []validation.ValidationError {
	errors := make([]validation.ValidationError, 0)
	err := utils.ValidateStruct(ctx, task)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		for _, err1 := range validateErr {
			errors = append(errors, err1.(validation.ValidationError))
		}
	}
	errMethod := utils.ValidateHttpMethod(task.Method)
	if errMethod != nil {
		errors = append(errors, errMethod)
	}
	return errors
}
