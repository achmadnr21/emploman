package utils

import (
	"fmt"
)

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %s", e.Message)
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict: %s", e.Message)
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not Found: %s", e.Message)
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("Internal Server Error: %s", e.Message)
}

type TooManyRequestError struct {
	Message string
}

func (e *TooManyRequestError) Error() string {
	return fmt.Sprintf("Too Many Request: %s", e.Message)
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad Request: %s", e.Message)
}

func GetHTTPErrorCode(err error) int {
	switch err.(type) {
	case *UnauthorizedError:
		return 401
	case *ConflictError:
		return 409
	case *NotFoundError:
		return 404
	case *InternalServerError:
		return 500
	case *TooManyRequestError:
		return 429
	case *BadRequestError:
		return 400
	default:
		return 500
	}
}
