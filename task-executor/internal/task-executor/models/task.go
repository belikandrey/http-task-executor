package models

const (
	// StatusNew - presents new task status.
	StatusNew = "new"
	// StatusError - presents error task status.
	StatusError = "error"
	// StatusInProcess - presents in_process task status.
	StatusInProcess = "in_process"
	// StatusDone - presents done task status.
	StatusDone = "done"
)

// Task presents task to execute.
// It contains ID, Status, ResponseStatus, ResponseLength, Headers.
type Task struct {
	// ID - unique identifier of task.
	ID int64 `db:"id"`
	// Status - current task status.
	Status string `db:"status"`
	// ResponseStatus - response status from the 3-rd service.
	ResponseStatus *int64 `db:"response_status_code"`
	// ResponseLength - response length from the 3-rd service.
	ResponseLength *int64 `db:"response_length"`
	// Headers - response headers from the 3-rd service.
	Headers []Header
}

// Header presents http header.
// It contains Name, Value, Input.
type Header struct {
	// Name - header name.
	Name string `db:"header_name"`
	// Value - header value.
	Value string `db:"header_value" validate:"required"`
	// Input - flag that points if header is from request (true) or response (false).
	Input bool `db:"header_input" validate:"required"`
}
