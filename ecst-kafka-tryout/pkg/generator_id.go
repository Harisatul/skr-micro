package pkg

import "github.com/google/uuid"

func GenerateId() (string, error) {
	v7, err := uuid.NewV7()
	return v7.String(), err
}
