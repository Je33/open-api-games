package domain

type TransactionType string

const (
	TransactionTypeDebit    TransactionType = "debit"
	TransactionTypeCredit   TransactionType = "credit"
	TransactionTypeRollback TransactionType = "rollback"
)

type Transaction struct {
	UID          string
	UserUID      string
	SessionUID   string
	Amount       int
	Currency     string
	Denomination int
	Type         TransactionType
}
