package response

import "fmt"

// BizError 业务错误。
type BizError struct {
	HTTPStatus int
	Code       int
	Message    string
	Err        error
}

func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(
			"%s: %v",
			e.Message,
			e.Err,
		)
	}

	return e.Message
}

func (e *BizError) Unwrap() error {
	return e.Err
}

// NewError 创建业务错误。
func NewError(
	httpStatus int,
	code int,
	message string,
) *BizError {
	return &BizError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// NewErrorWithCause 创建带原始错误的业务错误。
func NewErrorWithCause(
	httpStatus int,
	code int,
	message string,
	err error,
) *BizError {
	return &BizError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}
func BadRequest(message string) error {

	return NewError(
		400,
		CodeInvalidParams,
		message,
	)

}

func Unauthorized(message string) error {

	return NewError(
		401,
		CodeUnauthorized,
		message,
	)

}

func Forbidden(message string) error {

	return NewError(
		403,
		CodeForbidden,
		message,
	)

}

func NotFound(message string) error {

	return NewError(
		404,
		CodeNotFound,
		message,
	)

}

func InternalError(message string) error {

	return NewError(
		500,
		CodeInternalError,
		message,
	)

}
