package errorsHandler

import (
	"net/http"
)

const (
	UNAUTHENTICATEDError = "UNAUTHENTICATED_grpc"
	NotFoundUser         = "not_found_user"
	NameRequired         = "username_required"
	InternalError        = "internal_error"
	NotFoundResult       = "not_found_result"
)

var ErrorsMap = MapErrors{
	FA: FAErrorsMap,
	EN: ENErrorsMap,
}

var FAErrorsMap = LangErrorMap{
	UNAUTHENTICATEDError: "خطا در احراز هویت",
	NotFoundUser:         "کاربری یافت نشد",
	NameRequired:         "نام اجباری میباشد",
	InternalError:        "خطای داخلی",
}

var ENErrorsMap = LangErrorMap{
	UNAUTHENTICATEDError: "UNAUTHENTICATED",
	NotFoundUser:         "user not found",
	NameRequired:         "name is required filed",
	InternalError:        "Internal Error",
}

var CodeErrorsMap = CodeErrorMap{
	NotFoundUser: http.StatusNotFound,
	NameRequired: http.StatusBadRequest,
}
