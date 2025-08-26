package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"http-task-executor/task-service/internal/task-service/logger"
	"http-task-executor/task-service/internal/task-service/models"
	tasks2 "http-task-executor/task-service/internal/task-service/tasks"
	"http-task-executor/task-service/pkg/errors/general/validation"
	httpErrors "http-task-executor/task-service/pkg/errors/http"
	"http-task-executor/task-service/pkg/utils"
)

// TaskUseCase represents service layer to models.Task.
type TaskUseCase struct {
	log      logger.Logger
	repo     tasks2.Repository
	producer tasks2.Producer
}

// NewTaskUseCase creates new instance of TaskUseCase.
func NewTaskUseCase(log logger.Logger, repo tasks2.Repository, producer tasks2.Producer) *TaskUseCase {
	return &TaskUseCase{log: log, repo: repo, producer: producer}
}

// Create creates new task and send to producer.
func (t *TaskUseCase) Create(ctx context.Context, task *models.Task) (*models.Task, error) {

	validationErrors := validateTask(ctx, task)
	if len(validationErrors) > 0 {
		return nil, httpErrors.NewValidationError(validationErrors)
	}

	create, err := t.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	err = t.producer.Produce(create)
	if err != nil {
		errInternal := t.repo.Delete(ctx, create.ID)
		if errInternal != nil {
			return nil, httpErrors.NewInternalServerError(errInternal)
		}
		return nil, err
	}
	return create, nil
}

// GetByIdWithOutputHeaders returns models.Task by requested ID.
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

func validateTask(ctx context.Context, task *models.Task) []validation.TaskValidationError {
	errs := make([]validation.TaskValidationError, 0)
	err := utils.ValidateStruct(ctx, task)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		for _, err1 := range validateErr {
			errs = append(errs, err1.(validation.TaskValidationError))
		}
	}
	errMethod := utils.ValidateHttpMethod(task.Method)
	if errMethod != nil {
		errs = append(errs, errMethod)
	}
	return errs
}
