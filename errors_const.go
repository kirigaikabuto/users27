package users27

import "errors"

var (
	ErrNoUser = errors.New("no user by this username and password")
)
