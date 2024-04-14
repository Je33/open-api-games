package data

import "open-api-games/internal/domain"

func Currencies() []domain.Currency {
	return []domain.Currency{
		{
			Code:         "USD",
			Denomination: 2,
		},
	}
}
