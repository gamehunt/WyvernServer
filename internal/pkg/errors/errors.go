package errors

import "errors"

var AuthInvalidIdentity = errors.New("Invalid identity")
var UserAlreadyExists   = errors.New("User already exists")

var InvalidSession      = errors.New("Invalid session")
