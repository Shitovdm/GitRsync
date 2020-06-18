package helper

import "github.com/gofrs/uuid"

// GenerateUUID returns new UUID in v4 format.
func GenerateUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}