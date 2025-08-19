package models

const (
	StatusNew       = "new"
	StatusError     = "error"
	StatusInProcess = "in_process"
	StatusDone      = "done"
)

type Task struct {
	Id             int64  `db:"id"`
	Url            string `db:"url" validate:"required"`
	Method         string `db:"method" validate:"required"`
	Status         string `db:"status"`
	ResponseStatus *int64 `db:"response_status"`
	ResponseLength *int64 `db:"response_length"`
	Headers        []Header
}

type Header struct {
	Name  string `db:"header_name"`
	Value string `db:"header_value" validate:"required"`
	Input bool   `db:"header_input" validate:"required"`
}
