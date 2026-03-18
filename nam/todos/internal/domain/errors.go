package apperrors

// errors.go — Application Error Types
//
// Phase 1: Clean Architecture — Domain Layer
//
// This file is responsible for:
// 1. Define AppError struct with:
//    - Code    ErrorCode (enum: NotFound, InvalidParameter, AuthZ, AuthN, Internal)
//    - Message string
//    - Details map[string]string (optional metadata)
// 2. Define ErrorCode enum with constants
// 3. Define constructor functions:
//    - NewNotFound(message string) *AppError
//    - NewInvalidParameter(message string) *AppError
//    - NewAuthZ(message string) *AppError
//    - NewAuthN(message string) *AppError
//    - NewInternal(message string) *AppError
// 4. Implement the error interface: func (e *AppError) Error() string
//
// Why custom errors?
// - The handler layer maps AppError codes to gRPC status codes:
//   NotFound → codes.NotFound, AuthZ → codes.PermissionDenied, etc.
// - This keeps domain logic unaware of gRPC while still communicating error types
//
// See: resources/phase-01-architecture-grpc.md (error handling)
// See: resources/phase-01-architecture-grpc.md (AppError → gRPC status mapping)

type ErrorCode string

const (
	ErrorCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrorCodeInvalidParameter ErrorCode = "INVALID_PARAMETER"
	ErrorCodeAuthZ            ErrorCode = "AUTHZ"
	ErrorCodeAuthN            ErrorCode = "AUTHN"
	ErrorCodeInternal         ErrorCode = "INTERNAL"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]string
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Code:    ErrorCodeNotFound,
		Message: message,
	}
}
func NewInvalidParameter(message string) *AppError {
	return &AppError{
		Code:    ErrorCodeInvalidParameter,
		Message: message,
	}
}

func NewAuthZ(message string) *AppError {
	return &AppError{
		Code:    ErrorCodeAuthZ,
		Message: message,
	}
}

func NewAuthN(message string) *AppError {
	return &AppError{
		Code:    ErrorCodeAuthN,
		Message: message,
	}
}

func NewInternal(message string) *AppError {
	return &AppError{
		Code:    ErrorCodeInternal,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}
