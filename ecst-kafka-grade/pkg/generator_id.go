package pkg

import "github.com/google/uuid"

func GenerateId() (uuid.UUID, error) {
	v7, err := uuid.NewV7()
	return v7, err
}
