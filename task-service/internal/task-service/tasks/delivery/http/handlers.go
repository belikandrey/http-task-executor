package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"http-task-executor/task-service/internal/task-service/config"
	"http-task-executor/task-service/internal/task-service/logger"
	"http-task-executor/task-service/internal/task-service/tasks"
	"http-task-executor/task-service/internal/task-service/tasks/delivery/http/dto"
	"http-task-executor/task-service/internal/task-service/tasks/mapper"
	httpErrors "http-task-executor/task-service/pkg/errors/http"
)

// TaskHandlers represents handlers for router.
type TaskHandlers struct {
	cfg     *config.Config
	useCase tasks.UseCase
	logger  logger.Logger
}

// NewTaskHandlers creates a new TaskHandlers instance.
func NewTaskHandlers(
	cfg *config.Config,
	logger logger.Logger,
	useCase tasks.UseCase,
) *TaskHandlers {
	return &TaskHandlers{cfg: cfg, logger: logger, useCase: useCase}
}

// Create godoc
//
//	@Summary		Create task and execute request to 3rd service
//	@Description	Create task and execute request to 3rd service
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.NewTaskRequest	true	"Task create request"
//	@Success		201		{object}	dto.NewTaskResponse
//	@Router			/task [post].
func (h *TaskHandlers) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTaskRequest dto.NewTaskRequest

		err := render.DecodeJSON(r.Body, &newTaskRequest)
		if err != nil {
			h.logger.Error(err)
			code, data := httpErrors.ErrorResponse(err)
			render.Status(r, code)
			render.JSON(w, r, data)

			return
		}

		h.logger.Infof("Request body decoded %v", newTaskRequest)

		task := mapper.MapRequestToTask(&newTaskRequest)

		create, err := h.useCase.Create(r.Context(), &task)
		if err != nil {
			h.logger.Error(err)
			code, data := httpErrors.ErrorResponse(err)
			render.Status(r, code)
			render.JSON(w, r, data)

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, mapper.MapIDToTaskResponse(create.ID))
	}
}

// Get godoc
//
//	@Summary		Get task by id
//	@Description	Get task by id handler
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	dto.GetTaskResponse
//	@Router			/task/{id} [get].
func (h *TaskHandlers) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			h.logger.Error(err)
			code, data := httpErrors.ErrorResponse(err)
			render.Status(r, code)
			render.JSON(w, r, data)

			return
		}

		h.logger.Infof("Request path decoded %v", idInt)

		if idInt <= 0 {
			h.logger.Info("ID must be positive")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, httpErrors.NewRestError(http.StatusBadRequest, "Invalid id", nil))

			return
		}

		responseTask, err := h.useCase.GetByIDWithOutputHeaders(r.Context(), int64(idInt))
		if err != nil {
			h.logger.Error(err)
			code, data := httpErrors.ErrorResponse(err)
			render.Status(r, code)
			render.JSON(w, r, data)

			return
		}

		response := mapper.MapTaskToGetResponse(responseTask)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}
