package data

import "open-api-games/internal/domain"

func Sessions() []domain.Session {
	return []domain.Session{
		{
			UID:     "FIRST_SESSION_UID",
			UserUID: "FIRST_USER_UID",
		},
	}
}
