package postgres

import (
	"errors"
)

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violated foreign key")
	ErrViolatesUniqueKey  = errors.New("violates unique key")
	ErrUnknown            = errors.New("unknown")
)
