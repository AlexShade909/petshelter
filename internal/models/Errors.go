package models

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidReference = errors.New("shelter or clinic does not exist")
)
