package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

// 自定義的 通用 errors, 第四碼開頭為0
// ref: https://coda.io/d/_dXACODzm3DP/Golang-20201204_suA_M#_luWrx 新增修改時請通知並同步資訊
var (
	ErrInvalidInput = &_error{Code: "400001", Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}

	ErrUnauthorized = &_error{Code: "401001", Message: http.StatusText(http.StatusUnauthorized), Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}

	ErrResourceNotFound = &_error{Code: "404001", Message: "The specified resource does not exist.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}

	ErrMethodNotAllowed = &_error{Code: "405001", Message: "Server has received and recognized the request, but has rejected the specific HTTP method it’s using.", Status: http.StatusMethodNotAllowed, GRPCCode: codes.Unavailable}

	ErrResourceAlreadyExists = &_error{Code: "409004", Message: "The specified resource already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrResourceUnavailable   = &_error{Code: "409005", Message: "The specified resource is unavailable.", Status: http.StatusConflict, GRPCCode: codes.Unavailable}
	ErrResourceInsufficient  = &_error{Code: "409006", Message: "The specified resource is insufficient.", Status: http.StatusConflict, GRPCCode: codes.Unavailable}
	ErrInsufficientBalance   = &_error{Code: "409007", Message: "Insufficient balance", Status: http.StatusConflict, GRPCCode: codes.Unavailable}

	ErrInternalServerError = &_error{Code: "500000", Message: http.StatusText(http.StatusInternalServerError), Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
	ErrInternalError       = &_error{Code: "500001", Message: "The server encountered an internal error. Please retry the request.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
)

type _error struct {
	Status   int                    `json:"status"`
	Code     string                 `json:"code"`
	GRPCCode codes.Code             `json:"grpccode"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details"`
}

// HTTPError ...
type HTTPError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

func (e *_error) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(e.Code)
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

// Is ...
func (e *_error) Is(target error) bool {
	causeErr := errors.Cause(target)
	tErr, ok := causeErr.(*_error)
	if !ok {
		return false
	}
	return e.Code == tErr.Code
}

// NewWithMessage 抽換錯誤訊息
// 未定義的錯誤會被視為 ErrInternalError 類型
func NewWithMessage(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return WithStack(&_error{
			Status:   ErrInternalError.Status,
			Code:     ErrInternalError.Code,
			Message:  ErrInternalError.Message,
			GRPCCode: ErrInternalError.GRPCCode,
		})
	}
	err = &_error{
		Status:   _err.Status,
		Code:     _err.Code,
		Message:  fmt.Sprintf(message, args...),
		GRPCCode: _err.GRPCCode,
	}
	return WithStack(err)
}

// WithErrors 使用訂好的errors code 與訊息,如果未定義message 顯示對應的http status描述
func WithErrors(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return WithStack(&_error{
			Status:  ErrInternalError.Status,
			Code:    ErrInternalError.Code,
			Message: http.StatusText(ErrInternalError.Status),
		})
	}
	return WithStack(&_error{
		Status:  _err.Status,
		Code:    _err.Code,
		Message: _err.Message,
	})
}

// NewWithMessagef 抽換錯誤訊息
func NewWithMessagef(err error, format string, args ...interface{}) error {
	return NewWithMessage(err, fmt.Sprintf(format, args...))
}
