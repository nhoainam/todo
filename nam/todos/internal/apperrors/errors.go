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
	Code     ErrorCode
	Message  string
	metadata map[string]any
	cause    error
}

type MetadataOption func(map[string]any)

func WithMetadata(key string, value any) MetadataOption {
	return func(metadata map[string]any) {
		metadata[key] = value
	}
}

func NewNotFound(message string, cause error, metadata ...MetadataOption) *AppError {
	return &AppError{
		Code:    ErrorCodeNotFound,
		Message: message,
		cause:   cause,
		metadata: func() map[string]any {
			m := make(map[string]any)
			for _, option := range metadata {
				option(m)
			}
			return m
		}(),
	}
}
func NewInvalidParameter(message string, cause error, metadata ...MetadataOption) *AppError {
	return &AppError{
		Code:    ErrorCodeInvalidParameter,
		Message: message,
		cause:   cause,
		metadata: func() map[string]any {
			m := make(map[string]any)
			for _, option := range metadata {
				option(m)
			}
			return m
		}(),
	}
}

func NewAuthZ(message string, cause error, metadata ...MetadataOption) *AppError {
	return &AppError{
		Code:    ErrorCodeAuthZ,
		Message: message,
		cause:   cause,
		metadata: func() map[string]any {
			m := make(map[string]any)
			for _, option := range metadata {
				option(m)
			}
			return m
		}(),
	}
}

func NewAuthN(message string, cause error, metadata ...MetadataOption) *AppError {
	return &AppError{
		Code:    ErrorCodeAuthN,
		Message: message,
		cause:   cause,
		metadata: func() map[string]any {
			m := make(map[string]any)
			for _, option := range metadata {
				option(m)
			}
			return m
		}(),
	}
}

func NewInternal(message string, cause error, metadata ...MetadataOption) *AppError {
	return &AppError{
		Code:    ErrorCodeInternal,
		Message: message,
		cause:   cause,
		metadata: func() map[string]any {
			m := make(map[string]any)
			for _, option := range metadata {
				option(m)
			}
			return m
		}(),
	}
}

func (e *AppError) Error() string {
	return e.Message
}
