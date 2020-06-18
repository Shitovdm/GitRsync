package helper

import "github.com/gofrs/uuid"

func GenerateUuid() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}