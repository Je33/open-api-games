package data

import (
	"open-api-games/internal/domain"
)

func Users() []domain.User {
	return []domain.User{
		{
			UID:  "FIRST_USER_UID",
			Nick: "First User",
		},
	}
}
