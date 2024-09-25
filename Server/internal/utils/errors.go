package utils

import (
	"fmt"
	"log"
	"os"
)

type AppError struct {
	Code int
	Message string
	Err error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("AppError %d: %s - %v", e.Code, e.Message, e.Err)
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}