package validation

// TaskValidationError represents custom error on validation fields.
type TaskValidationError interface {
	Error() string
	Field() string
	ActualTag() string
}

// CustomFiledError represents implementation of TaskValidationError for using in fields validation.
type CustomFiledError struct {
	Fld string `json:"field"`
	Msg string `json:"message"`
	Tag string `json:"actualTag"`
}

// Error returns error.
func (err CustomFiledError) Error() string {
	return err.Msg
}

// Field returns error field name.
func (err CustomFiledError) Field() string {
	return err.Fld
}

// ActualTag returns tag on error field.
func (err CustomFiledError) ActualTag() string {
	return err.Tag
}
