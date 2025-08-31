package domain

import "errors"

var (
	ErrValidation                = errors.New("validation error")
	ErrEntrepreneurNotFound      = errors.New("entrepreneur not found")
	ErrEntrepreneurAlreadyExists = errors.New("entrepreneur already exists")
)
