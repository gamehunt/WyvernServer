package repository

import (
	"wyvern/server/internal/types"
)

type RegistrationCache interface {
	Save(userId types.ID) (*types.ID, error)
	Restore(registrationId types.ID) (*types.ID, error)
}
