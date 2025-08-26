package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"http-task-executor/task-service/internal/task-service/models"
	"http-task-executor/task-service/internal/task-service/tasks/delivery/http/dto"
	"http-task-executor/task-service/internal/task-service/tasks/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandlers_Create(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	input := fmt.Sprintf(`{"url": "%s", "method": "%s"}`, "http://test.com", "GET")

	request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(input)))
	request.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()

	resTask := &models.Task{Url: "http://test.com", Method: "GET", Id: int64(1)}

	mockUseCase.EXPECT().Create(context.Background(), gomock.Any()).Return(resTask, nil)

	handlers.Create().ServeHTTP(res, request)

	body := res.Body.String()

	var response dto.NewTaskResponse

	require.Equal(t, http.StatusOK, res.Code)

	require.NoError(t, json.Unmarshal([]byte(body), &response))

	require.Equal(t, resTask.Id, response.Id)
}

func TestTaskHandlers_CreateWithHeaders(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	input := fmt.Sprintf(`{"url": "%s", "method": "%s", "headers" : {"%s":"%s"}}`, "http://test.com", "GET", "TEST1", "TEST1_VALUE")

	request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(input)))
	request.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()

	resTask := &models.Task{Url: "http://test.com", Method: "GET", Id: int64(1)}

	mockUseCase.EXPECT().Create(context.Background(), gomock.Any()).Return(resTask, nil)

	handlers.Create().ServeHTTP(res, request)

	body := res.Body.String()

	var response dto.NewTaskResponse

	require.Equal(t, http.StatusOK, res.Code)

	require.NoError(t, json.Unmarshal([]byte(body), &response))

	require.Equal(t, resTask.Id, response.Id)
}

func TestTaskHandlers_CreateWithErrorInUC(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	input := fmt.Sprintf(`{"url": "%s", "method": "%s"}`, "http://test.com", "GET")

	request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(input)))
	request.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()

	mockUseCase.EXPECT().Create(context.Background(), gomock.Any()).Return(nil, errors.New("error"))

	handlers.Create().ServeHTTP(res, request)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	body := res.Body.String()

	require.NotEmpty(t, body)
}

func TestTaskHandlers_CreateWithInvalidJSONSyntax(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	input := fmt.Sprintf(`url": "%s", "method": "%s"}`, "http://test.com", "GET")

	request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(input)))
	request.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handlers.Create().ServeHTTP(res, request)

	require.Equal(t, http.StatusBadRequest, res.Code)
	body := res.Body.String()

	require.NotEmpty(t, body)
}

func TestTaskHandlers_CreateWithUnmarshalTypeError(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	input := fmt.Sprintf(`"url": "%s", "method": "%s"}`, "http://test.com", "GET")

	request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(input)))
	request.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handlers.Create().ServeHTTP(res, request)

	require.Equal(t, http.StatusBadRequest, res.Code)
	body := res.Body.String()

	require.NotEmpty(t, body)
}

func TestTaskHandlers_Get(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	request := httptest.NewRequest(http.MethodGet, "/task/{id}", nil)
	request.Header.Add("Content-Type", "application/json")

	params := make(map[string]string)
	params["id"] = "1"
	request = addChiURLParams(request, params)

	res := httptest.NewRecorder()

	respStatus := int64(200)
	respLength := int64(10)

	resTask := &models.Task{Id: 1, Status: models.StatusInProcess, ResponseStatus: &respStatus, ResponseLength: &respLength}

	mockUseCase.EXPECT().GetByIdWithOutputHeaders(gomock.Any(), gomock.Any()).Return(resTask, nil)

	handlers.Get().ServeHTTP(res, request)

	body := res.Body.String()

	var response dto.GetTaskResponse

	require.Equal(t, http.StatusOK, res.Code)

	require.NoError(t, json.Unmarshal([]byte(body), &response))

	require.Equal(t, resTask.Id, response.ID)
	require.Equal(t, resTask.Status, response.Status)
	require.Equal(t, resTask.ResponseStatus, response.ResponseStatus)
	require.Equal(t, resTask.ResponseLength, response.ResponseLength)
}

func TestTaskHandlers_GetStringId(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	request := httptest.NewRequest(http.MethodGet, "/task/{id}", nil)
	request.Header.Add("Content-Type", "application/json")

	params := make(map[string]string)
	params["id"] = "asdads"
	request = addChiURLParams(request, params)

	res := httptest.NewRecorder()

	handlers.Get().ServeHTTP(res, request)

	body := res.Body

	require.Equal(t, http.StatusBadRequest, res.Code)

	require.NotEmpty(t, body)
}

func TestTaskHandlers_GetNegativeId(t *testing.T) {
	t.Parallel()
	ctrx := gomock.NewController(t)
	defer ctrx.Finish()

	sugar := zap.New(zapcore.NewNopCore()).Sugar()

	mockUseCase := mock.NewMockUseCase(ctrx)

	handlers := NewTaskHandlers(nil, sugar, mockUseCase)

	request := httptest.NewRequest(http.MethodGet, "/task/{id}", nil)
	request.Header.Add("Content-Type", "application/json")

	params := make(map[string]string)
	params["id"] = "-15"
	request = addChiURLParams(request, params)

	res := httptest.NewRecorder()

	handlers.Get().ServeHTTP(res, request)

	body := res.Body

	require.Equal(t, http.StatusBadRequest, res.Code)

	require.NotEmpty(t, body)
}

func addChiURLParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
