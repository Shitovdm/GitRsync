package helper

import "github.com/gofrs/uuid"

func GenerateUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}