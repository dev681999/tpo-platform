package utils

import "github.com/google/uuid"

// ParseUUID parse UUID and returns any error
func ParseUUID(id string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return uuid, ErrUUIDParse
	}

	return uuid, nil
}
