package domain

type ProcessBalanceReq struct {
	GameSessionUID string
	Currency       string
}

type ProcessBalanceRes struct {
	UserUID      string
	UserNick     string
	Amount       int
	Currency     string
	Denomination int
	MaxWin       int
	JpKey        string
}

type ProcessDebitCreditRollbackReq struct {
	TransactionUID string
	GameSessionUID string
	UserUID        string
	UserNick       string
	Amount         int
	Currency       string
	Denomination   int
	MaxWin         int
	JpKey          string
	SpinMeta       string
	BetMeta        string
}

type ProcessDebitCreditRollbackRes struct {
	TransactionUID string
	UserNick       string
	Amount         int
	Currency       string
	Denomination   int
	MaxWin         int
}

type ProcessApiDataApi string

const (
	ProcessApiDataApiRoundComplete ProcessApiDataApi = "roundComplete"
)

type ProcessApiDataData struct {
	BetUID string
}

type ProcessMetaDataReq struct {
	UserUID        string
	GameSessionUID string
	Currency       string
	Api            ProcessApiDataApi
	Data           ProcessApiDataData
}

type ProcessMetaDataRes struct {
	Api  ProcessApiDataApi
	Data string
}
