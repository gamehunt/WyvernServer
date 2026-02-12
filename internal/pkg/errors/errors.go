package errors

import "errors"

var AuthInvalidIdentity = errors.New("Invalid identity")
var AuthInvalidCredId   = errors.New("Invalid cred id")
var UserAlreadyExists   = errors.New("User already exists")

var InvalidSession      = errors.New("Invalid session")
