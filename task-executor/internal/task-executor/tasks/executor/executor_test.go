package executor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"http-task-executor/task-executor/internal/task-executor/dto"
	"http-task-executor/task-executor/internal/task-executor/models"
	"http-task-executor/task-executor/internal/task-executor/tasks/mock"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const duration = 3 * time.Second

func TestExecutor_Execute(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)

	closer := io.NopCloser(strings.NewReader(`{"message": "ok"}`))
	mockResp := &http.Response{
		StatusCode:    200,
		Body:          io.NopCloser(strings.NewReader(`{"message": "ok"}`)),
		Header:        make(http.Header),
		ContentLength: int64(len(`{"message": "ok"}`)),
	}

	written, err := io.Copy(io.Discard, closer)

	mockTransport := &mockRoundTripper{
		Response: mockResp,
		Err:      nil,
	}

	provider := newMockClientProvider(mockTransport)

	executor := NewExecutor(sugar, mockTasksRepo, provider, duration)

	message := dto.KafkaTaskMessage{
		ID:     1,
		Method: "GET",
		Url:    "https://www.google.com",
	}

	marshal, err := json.Marshal(message)
	require.NoError(t, err)

	_, cancel := context.WithTimeout(context.Background(), duration)

	mockTasksRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	status := int64(200)

	require.NoError(t, err)

	expectedTaskToUpdate := models.Task{
		Id:             message.ID,
		Status:         models.StatusDone,
		ResponseStatus: &status,
		ResponseLength: &written,
		Headers:        make([]models.Header, 0),
	}

	mockTasksRepo.EXPECT().UpdateResult(gomock.Any(), gomock.Cond(func(x any) bool {
		task := x.(*models.Task)
		return task.Id == expectedTaskToUpdate.Id &&
			task.Status == expectedTaskToUpdate.Status &&
			*task.ResponseStatus == *expectedTaskToUpdate.ResponseStatus &&
			*task.ResponseLength == *expectedTaskToUpdate.ResponseLength &&
			len(task.Headers) == len(expectedTaskToUpdate.Headers)
	},
	)).Times(1)

	executor.Execute(marshal)

	defer cancel()
}

func TestExecutor_ExecuteTaskWithHeader(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)

	headerName := "TEST"
	headerValue := "TEST_VALUE"
	closer := io.NopCloser(strings.NewReader(`{"message": "ok"}`))
	header := make(http.Header)
	header[headerName] = []string{headerValue}
	mockResp := &http.Response{
		StatusCode:    200,
		Body:          io.NopCloser(strings.NewReader(`{"message": "ok"}`)),
		Header:        header,
		ContentLength: int64(len(`{"message": "ok"}`)),
	}

	written, err := io.Copy(io.Discard, closer)

	mockTransport := &mockRoundTripper{
		Response: mockResp,
		Err:      nil,
	}

	provider := newMockClientProvider(mockTransport)

	executor := NewExecutor(sugar, mockTasksRepo, provider, duration)

	message := dto.KafkaTaskMessage{
		ID:     1,
		Method: "GET",
		Url:    "https://www.google.com",
	}

	marshal, err := json.Marshal(message)
	require.NoError(t, err)

	_, cancel := context.WithTimeout(context.Background(), duration)

	mockTasksRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	status := int64(200)

	require.NoError(t, err)

	headers := make([]models.Header, 0)
	headers = append(headers, models.Header{Name: headerName, Value: headerValue, Input: false})
	expectedTaskToUpdate := models.Task{
		Id:             message.ID,
		Status:         models.StatusDone,
		ResponseStatus: &status,
		ResponseLength: &written,
		Headers:        headers,
	}

	mockTasksRepo.EXPECT().UpdateResult(gomock.Any(), gomock.Cond(func(x any) bool {
		task := x.(*models.Task)
		return task.Id == expectedTaskToUpdate.Id &&
			task.Status == expectedTaskToUpdate.Status &&
			*task.ResponseStatus == *expectedTaskToUpdate.ResponseStatus &&
			*task.ResponseLength == *expectedTaskToUpdate.ResponseLength &&
			len(task.Headers) == len(expectedTaskToUpdate.Headers) &&
			task.Headers[0].Name == expectedTaskToUpdate.Headers[0].Name &&
			task.Headers[0].Value == expectedTaskToUpdate.Headers[0].Value &&
			task.Headers[0].Input == expectedTaskToUpdate.Headers[0].Input
	},
	)).Times(1)

	executor.Execute(marshal)

	defer cancel()
}

func TestExecutor_ExecuteTaskWithNotWoringUpdateStatus(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockTasksRepo := mock.NewMockRepository(ctrx)

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"message": "ok"}`)),
		Header:     make(http.Header),
	}

	mockTransport := &mockRoundTripper{
		Response: mockResp,
		Err:      nil,
	}

	provider := newMockClientProvider(mockTransport)

	executor := NewExecutor(sugar, mockTasksRepo, provider, duration)

	message := dto.KafkaTaskMessage{
		ID:     1,
		Method: "GET",
		Url:    "https://www.google.com",
	}

	marshal, err := json.Marshal(message)
	require.NoError(t, err)

	_, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	mockTasksRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

	mockTasksRepo.EXPECT().UpdateResult(gomock.Any(), gomock.Any()).Times(0)

	executor.Execute(marshal)

}

type mockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

type mockClientProvider struct {
	transport *mockRoundTripper
}

func newMockClientProvider(transport *mockRoundTripper) *mockClientProvider {
	return &mockClientProvider{transport: transport}
}

func (c *mockClientProvider) Client() *http.Client {
	return &http.Client{Transport: c.transport}
}
