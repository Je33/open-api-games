package domain

import "github.com/google/uuid"

func GenUID() string {
	uid := uuid.New()
	return uid.String()
}
