package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"http-task-executor/internal/models"
	"http-task-executor/internal/tasks/mock"
	errorsHttp "http-task-executor/pkg/errors/http"
	"net/http"
	"testing"
	"time"
)

func TestTaskUseCase_Create(t *testing.T) {
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)
	mockExecutor := mock.NewMockExecutor(ctrx)

	useCase := NewTaskUseCase(sugar, mockTasksRepo, mockExecutor)

	task := &models.Task{
		Method: "GET",
		Url:    "https://www.google.com",
		Status: models.StatusNew,
	}

	ctx := context.Background()

	mockTasksRepo.EXPECT().Create(ctx, gomock.Eq(task)).Return(task, nil).Times(1)
	called := make(chan struct{}, 1)
	mockExecutor.EXPECT().ExecuteTask(gomock.Any()).Do(func(task models.Task) {
		called <- struct{}{}
	}).Times(1)

	create, err := useCase.Create(ctx, task)

	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, create)

	select {
	case <-called:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Expected ExecuteTask to be called in goroutine")
	}
}

func TestTaskUseCase_CreateWithErrorsNotExecuteTask(t *testing.T) {
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)
	mockExecutor := mock.NewMockExecutor(ctrx)

	useCase := NewTaskUseCase(sugar, mockTasksRepo, mockExecutor)

	task := &models.Task{
		Method: "GET",
		Url:    "https://www.google.com",
		Status: models.StatusNew,
	}

	ctx := context.Background()

	mockTasksRepo.EXPECT().Create(ctx, gomock.Eq(task)).Return(nil, errors.New("error"))
	called := make(chan struct{}, 1)
	mockExecutor.EXPECT().ExecuteTask(gomock.Any()).Do(func(task models.Task) {
		called <- struct{}{}
	}).Times(0)

	create, err := useCase.Create(ctx, task)

	require.Error(t, err)
	require.Nil(t, create)

	select {
	case <-called:
		t.Fatal("Expected ExecuteTask to NOT called in goroutine")
	case <-time.After(500 * time.Millisecond):
	}
}

func TestTaskUseCase_CreateWithInvalidMethodNotExecuteTask(t *testing.T) {
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)
	mockExecutor := mock.NewMockExecutor(ctrx)

	useCase := NewTaskUseCase(sugar, mockTasksRepo, mockExecutor)

	task := &models.Task{
		Method: "tersfasd",
		Url:    "https://www.google.com",
		Status: models.StatusNew,
	}

	ctx := context.Background()

	mockTasksRepo.EXPECT().Create(ctx, gomock.Eq(task)).Times(0)
	called := make(chan struct{}, 1)
	mockExecutor.EXPECT().ExecuteTask(gomock.Any()).Do(func(task models.Task) {
		called <- struct{}{}
	}).Times(0)

	create, err := useCase.Create(ctx, task)

	require.Error(t, err)
	require.Nil(t, create)
	require.NotEmpty(t, err.(errorsHttp.RestError))
	require.Equal(t, err.(errorsHttp.RestError).ErrStatus, http.StatusBadRequest)

	select {
	case <-called:
		t.Fatal("Expected ExecuteTask to NOT called in goroutine")
	case <-time.After(500 * time.Millisecond):
	}
}

func TestTaskUseCase_CreateWithInvalidUrlNotExecuteTask(t *testing.T) {
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)
	mockExecutor := mock.NewMockExecutor(ctrx)

	useCase := NewTaskUseCase(sugar, mockTasksRepo, mockExecutor)

	task := &models.Task{
		Method: "GET",
		Url:    ":/www.goog",
		Status: models.StatusNew,
	}

	ctx := context.Background()

	mockTasksRepo.EXPECT().Create(ctx, gomock.Eq(task)).Times(0)
	called := make(chan struct{}, 1)
	mockExecutor.EXPECT().ExecuteTask(gomock.Any()).Do(func(task models.Task) {
		called <- struct{}{}
	}).Times(0)

	create, err := useCase.Create(ctx, task)

	require.Error(t, err)
	require.Nil(t, create)
	require.NotEmpty(t, err.(errorsHttp.RestError))
	require.Equal(t, err.(errorsHttp.RestError).ErrStatus, http.StatusBadRequest)

	select {
	case <-called:
		t.Fatal("Expected ExecuteTask to NOT called in goroutine")
	case <-time.After(500 * time.Millisecond):
	}
}

func TestTaskUseCase_GetByIdWithOutputHeadersInvalidId(t *testing.T) {
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)
	mockExecutor := mock.NewMockExecutor(ctrx)

	useCase := NewTaskUseCase(sugar, mockTasksRepo, mockExecutor)

	id := int64(-1)

	ctx := context.Background()

	mockTasksRepo.EXPECT().GetByIdWithOutputHeaders(ctx, id).Times(0)

	task, err := useCase.GetByIdWithOutputHeaders(ctx, id)

	require.Error(t, err)
	require.Nil(t, task)
	require.NotEmpty(t, err.(errorsHttp.RestError))
	require.Equal(t, err.(errorsHttp.RestError).ErrStatus, http.StatusBadRequest)

}

func TestTaskUseCase_GetByIdWithOutputHeadersValidId(t *testing.T) {
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)
	mockExecutor := mock.NewMockExecutor(ctrx)

	useCase := NewTaskUseCase(sugar, mockTasksRepo, mockExecutor)

	id := int64(15)

	task := &models.Task{
		Id:     id,
		Method: "GET",
		Url:    "https://www.google.com",
		Status: models.StatusNew,
	}

	ctx := context.Background()

	mockTasksRepo.EXPECT().GetByIdWithOutputHeaders(ctx, id).Return(task, nil).Times(1)

	returnedTask, err := useCase.GetByIdWithOutputHeaders(ctx, id)

	require.NoError(t, err)
	require.NotNil(t, task)
	assert.Equal(t, task.Id, returnedTask.Id)
	assert.Equal(t, task.Status, returnedTask.Status)
	assert.Equal(t, task.Method, returnedTask.Method)
	assert.Equal(t, task.Url, returnedTask.Url)
}
