package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUnsupportedCity        = errors.New("unsupported city")
	ErrReceptionAlreadyExists = errors.New("reception aldeady open")
	ErrNoOpenReception        = errors.New("no open reception for pvz")
	ErrNoProducts             = errors.New("no found any product")
	ErrNoPVZ                  = errors.New("no found any PVZ")
)
