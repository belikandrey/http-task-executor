package validation

type ValidationError interface {
	Error() string
	Field() string
	ActualTag() string
}

type CustomFiledError struct {
	Fld string `json:"field"`
	Msg string `json:"message"`
	Tag string `json:"actualTag"`
}

func (err CustomFiledError) Error() string {
	return err.Msg
}

func (err CustomFiledError) Field() string {
	return err.Fld
}

func (err CustomFiledError) ActualTag() string {
	return err.Tag
}
