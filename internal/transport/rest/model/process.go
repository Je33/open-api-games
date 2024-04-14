package model

type ProcessApiCommand string

const (
	ProcessApiCommandBalance  ProcessApiCommand = "balance"
	ProcessApiCommandDebit    ProcessApiCommand = "debit"
	ProcessApiCommandCredit   ProcessApiCommand = "credit"
	ProcessApiCommandRollback ProcessApiCommand = "rollback"
	ProcessApiCommandMetaData ProcessApiCommand = "metaData"
)

func (p ProcessApiCommand) String() string {
	return string(p)
}

func (p ProcessApiCommand) IsValid() bool {
	return p == ProcessApiCommandBalance || p == ProcessApiCommandDebit || p == ProcessApiCommandCredit || p == ProcessApiCommandRollback || p == ProcessApiCommandMetaData
}

type ProcessApiReqData interface {
	*ProcessBalanceReq | *ProcessDebitCreditRollbackReq | *ProcessMetaDataReq
}

type ProcessApiResData interface {
	*ProcessBalanceRes | *ProcessDebitCreditRollbackRes | *ProcessMetaDataRes
}

type ProcessCommand struct {
	Api ProcessApiCommand `json:"api"`
}

type ProcessReq[T ProcessApiReqData] struct {
	Api  ProcessApiCommand `json:"api"`
	Data T                 `json:"data"`
}

type ProcessRes[T ProcessApiResData] struct {
	Api       ProcessApiCommand `json:"api"`
	Data      T                 `json:"data"`
	IsSuccess bool              `json:"isSuccess"`
	Error     string            `json:"error"`
	ErrorMsg  string            `json:"errorMsg"`
}

type ProcessBalanceReq struct {
	GameSessionUID string `json:"gameSessionId"`
	Currency       string `json:"currency"`
}

type ProcessDebitCreditRollbackReq struct {
	TransactionUID string `json:"transactionId"`
	GameSessionUID string `json:"gameSessionId"`
	UserUID        string `json:"userId"`
	UserNick       string `json:"userNick"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	Denomination   int    `json:"denomination"`
	MaxWin         int    `json:"maxWin"`
	JpKey          string `json:"jpKey"`
	SpinMeta       string `json:"spinMeta"`
	BetMeta        string `json:"betMeta"`
}

type ProcessMetaDataReq struct {
	UserUID        string             `json:"userId"`
	GameSessionUID string             `json:"gameSessionId"`
	Currency       string             `json:"currency"`
	Api            ProcessApiDataApi  `json:"api"`
	Data           ProcessApiDataData `json:"data"`
}

type ProcessApiDataApi string

const (
	ProcessApiDataApiRoundComplete ProcessApiDataApi = "roundComplete"
)

type ProcessApiDataData struct {
	BetId string `json:"betId"`
}

type ProcessBalanceRes struct {
	UserUID      string `json:"userId"`
	UserNick     string `json:"userNick"`
	Amount       int    `json:"amount"`
	Currency     string `json:"currency"`
	Denomination int    `json:"denomination"`
	MaxWin       int    `json:"maxWin"`
	JpKey        string `json:"jpKey"`
}

type ProcessDebitCreditRollbackRes struct {
	TransactionUID string `json:"transactionId"`
	UserNick       string `json:"userNick"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	Denomination   int    `json:"denomination"`
	MaxWin         int    `json:"maxWin"`
}

type ProcessMetaDataRes struct {
	Api  ProcessApiDataApi `json:"api"`
	Data string            `json:"data"`
}
