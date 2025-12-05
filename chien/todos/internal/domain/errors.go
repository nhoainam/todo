package domain

// errors.go — Application Error Types
//
// Week 1: Clean Architecture — Domain Layer
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
// See: resources/week-01-clean-architecture.md (error handling)
// See: resources/week-02-grpc-protobuf.md (AppError → gRPC status mapping)
