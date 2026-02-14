package errors

import "errors"

var AuthInvalidIdentity = errors.New("Invalid identity")
var UserAlreadyExists   = errors.New("User already exists")
var InvalidSession      = errors.New("Invalid session")

var InvalidMember       = errors.New("No such member")
var InvalidGuild        = errors.New("No such guild")
var InvalidChannel      = errors.New("No such channel")
var InvalidMessage      = errors.New("No such message")
var InvalidRole         = errors.New("No such role")
var InsufficientPerms   = errors.New("Insufficient permissions")
