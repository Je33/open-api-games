package data

import "open-api-games/internal/domain"

func Balances() []domain.Balance {
	return []domain.Balance{
		{
			UserUID:      "FIRST_USER_UID",
			Amount:       1000,
			Currency:     "USD",
			Denomination: 2,
		},
	}
}
