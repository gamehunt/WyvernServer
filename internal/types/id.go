package types

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() (ID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return ID(uuid.Nil), err
	}
	return ID(id), nil
}

func ParseID(s string) (ID, error) {
	u, err := uuid.Parse(s)
	return ID(u), err
}

func (id ID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid UUID string: %w", err)
	}

	*id = ID(u)
	return nil
}
