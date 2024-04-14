package domain

import (
	"fmt"
	"strings"
)

const (
	externalCode   = "external"
	externalSource = "[external]"
)

type Error struct {
	Code     string  `json:"code"`
	Source   string  `json:"source"`
	Internal []Error `json:"internal"`
	External []error `json:"external"`
}

func NewError(source string) *Error {
	return &Error{Source: source}
}

func (e Error) Error() string {
	err := e.String()
	if len(e.External) > 0 {
		err = fmt.Sprintf("%s: %s", err, e.GetExternal())
	}
	return fmt.Sprintf("%s: %s", err, e.GetInternal())
}

func (e *Error) GetInternal() string {
	var out []string
	for _, err := range e.Internal {
		out = append(out, err.GetInternal())
		if len(err.External) > 0 {
			out = append(out, err.GetExternal())
		}
	}
	return strings.Join(out, " / ")
}

func (e *Error) GetExternal() string {
	var out []string
	for _, nErr := range e.External {
		out = append(out, nErr.Error())
	}
	return strings.Join(out, " + ")
}

func (e *Error) String() string {
	return fmt.Sprintf("%s:%s", e.Source, e.Code)
}

func (e *Error) SetCode(code string) *Error {
	e.Code = code
	return e
}

func (e *Error) SetSource(source string) *Error {
	e.Source = source
	return e
}

func (e *Error) Add(err error) *Error {
	if err == nil {
		return e
	}
	cErr, ok1 := err.(*Error)
	dErr, ok2 := err.(Error)
	if ok1 {
		e.Internal = append(e.Internal, *cErr)
	} else if ok2 {
		e.Internal = append(e.Internal, dErr)
	} else {
		e.External = append(e.External, err)
	}
	return e
}

func AsError(err error) *Error {
	cErr, ok1 := err.(*Error)
	if ok1 {
		return cErr
	}
	dErr, ok2 := err.(Error)
	if ok2 {
		return &dErr
	}
	err2 := NewError(externalSource).SetCode(externalCode).Add(err)
	return err2
}

const (
	ErrNone                   = "NO_ERROR"
	ErrServer                 = "INTERNAL_ERROR"
	ErrConfig                 = "CONFIG_ERROR"
	ErrConnect                = "CONNECTION_ERROR"
	ErrRepoInit               = "REPO_INIT_ERROR"
	ErrInvalidApiCommand      = "INVALID_API_COMMAND"
	ErrInvalidRequest         = "INVALID_REQUEST"
	ErrProcessingRequest      = "PROCESSING_REQUEST_ERROR"
	ErrNotFound               = "NOT_FOUND"
	ErrDecrement              = "INSUFFICIENT_BALANCE"
	ErrIncrement              = "CREDIT_ERROR"
	ErrInvalidTransactionType = "TRANSACTION_TYPE_INVALID"
	ErrRollback               = "ROLLBACK_ERROR"
	ErrSessionNotFound        = "SESSION_NOT_FOUND"
	ErrUserNotFound           = "USER_NOT_FOUND"
	ErrBalanceNotFound        = "BALANCE_NOT_FOUND"
	ErrRepoCreate             = "REPO_CREATE_ERROR"
	ErrEmptyTransactionUID    = "EMPTY_TRANSACTION_ID"
	ErrTransactionNotFound    = "TRANSACTION_NOT_FOUND"
	ErrUnknownCurrency        = "UNKNOWN_CURRENCY"
	ErrSignInvalid            = "INVALID_SIGN"
	ErrSignEmpty              = "SIGN_NOT_PROVIDED"
	ErrReadBody               = "READ_BODY_ERROR"
)
