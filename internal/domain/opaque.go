package domain

import "wyvern/server/internal/types"

type OpaqueLoginSession struct {
	UserId types.ID `json:"user_id"`
	Data   []byte   `json:"data"`
}
