package db

import "errors"

var ErrRowAlreadyExists = errors.New("row already exists")
var ErrNotFound = errors.New("not found")
