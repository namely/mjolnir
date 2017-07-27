package interceptor

import (
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ErrorKeyer declares something that fulfills error interface and has a key.
// We will use this when we want machine readable error string alongside error message (e.g. proto errors)
// implements error interface in order to be passed around as an error
type ErrorKeyer interface {
	error
	Key() string
}

// ErrorFielder declares something that fulfills error interface and has Logrus error fields.
// We will use this to pass error fields from handlers to middlewares for Logs
// implements error interface in order to be passed around as an error
type ErrorFielder interface {
	error
	Fields() *logrus.Fields
}

// KeyedErr implements ErrorKeyer
type KeyedErr struct {
	error
	Message  string
	ErrorKey string
}

// FieldErr implements ErrorFielder
type FieldErr struct {
	error
	ErrFields *logrus.Fields
}

// NewFieldErr creates a new error with fields
func NewFieldErr(fields *logrus.Fields) *FieldErr {
	return &FieldErr{ErrFields: fields}
}

// Fields returns the fields for this error
func (e *FieldErr) Fields() *logrus.Fields {
	return e.ErrFields
}

// Error returns the error message
// Needed to fulfill interface requirement (error)
func (e *FieldErr) Error() string {
	return "Don't use me, just use my Fields() method for logs!"
}

// Error returns the error message
func (e *KeyedErr) Error() string {
	return e.Message
}

// Key returns the error key
func (e *KeyedErr) Key() string {
	return e.ErrorKey
}

// ErrGrpcInternalError indicates an internal server error
var ErrGrpcInternalError = grpc.Errorf(codes.Unknown, "internal server error")
