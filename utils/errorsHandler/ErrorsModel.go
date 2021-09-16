package errorsHandler

import "net/http"

const (
	DefaultError = "default_error"
)

type MapErrors map[string]LangErrorMap
type LangErrorMap map[string]string
type CodeErrorMap map[string]int

type ErrorsModel struct {
	ErrorsMap     MapErrors
	CodeErrorsMap CodeErrorMap
}

var DefaultErrorsMap = MapErrors{
	FA: DefaultFAErrorsMap,
	EN: DefaultENErrorsMap,
}

var DefaultFAErrorsMap = LangErrorMap{
	DefaultError: "مشکل در ارتباط داخلی",
}

var DefaultENErrorsMap = LangErrorMap{
	DefaultError: "internal error",
}

var DefaultCodeErrorsMap = CodeErrorMap{
	DefaultError: http.StatusInternalServerError,
}
