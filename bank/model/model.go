package model

type Account struct {
	ID              string `json:"account_id"`
	Ledger          uint32 `json:"ledger"`
	Code            uint16 `json:"code"`
	IsLinked        bool   `json:"is_linked"`
	IsDebitBalance  bool   `json:"is_debit_balance"`
	IsCreditBalance bool   `json:"is_credit_balance"`
}

type Balance struct {
	AccountID        string `json:"account_id"`
	AvailableBalance uint64 `json:"availabe_balance"`
	AccountBalance   uint64 `json:"account_balance"`
}

type Transfer struct {
	ID              string `json:"trasnfer_id"`
	DebitAccountID  string `json:"debit_account_id"`
	CreditAccountID string `json:"credit_account_id"`
	PendingID       string `json:"pending_id"`
	Ledger          uint32 `json:"ledger"`
	Code            uint16 `json:"code"`
	Amount          uint64 `json:"amount"`
	IsLinked        bool   `json:"is_linked"`
	IsPending       bool   `json:"is_pending_trasnfer"`
	IsPostPending   bool   `json:"is_post_pending_id"`
	IsVoidPending   bool   `json:"is_void_pending_id"`
}

type Authorization struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
	Timestamp uint64
	PendingID string
}

type Presentment struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
}

var SignalChannels = struct {
	AUTHORIZATION_CHANNEL string
	PRESENTMENT_CHANNEL   string
}{
	AUTHORIZATION_CHANNEL: "AUTHORIZATION_CHANNEL",
	PRESENTMENT_CHANNEL:   "PRESENTMENT_CHANNEL",
}
