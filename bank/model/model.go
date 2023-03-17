package model

type Account struct {
	ID              string `json:"account_id"`
	Ledger          uint32 `json:"ledger"`
	Code            uint16 `json:"code"`
	IsLinked        bool   `json:"is_linked"`
	IsDebitBalance  bool   `json:"is_debit_balance"`
	IsCreditBalance bool   `json:"is_credit_balance"`
	DebitsPending   uint64 `json:"debits_pending"`
	DebitsPosted    uint64 `json:"debits_posted"`
	CreditsPending  uint64 `json:"credits_pending"`
	CreditsPosted   uint64 `json:"credits_posted"`
}

type Balance struct {
	AccountID        string `json:"account_id"`
	AvailableBalance uint64 `json:"availabe_balance"`
	AccountBalance   uint64 `json:"account_balance"`
}

type Transfer struct {
	DebitAccountID  string `json:"debit_account_id"`
	CreditAccountID string `json:"credit_account_id"`
	PendingID       string `json:"pending_id" default:"0"`
	Amount          uint64 `json:"amount"`
	IsLinked        bool   `json:"is_linked" default:"false"`
	IsPending       bool   `json:"is_pending_trasnfer" default:"false"`
	IsPostPending   bool   `json:"is_post_pending_id" default:"false"`
	IsVoidPending   bool   `json:"is_void_pending_id" default:"false"`
}

type Authorization struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
	PendingID string
}

type Presentment struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
	PendingID string
}

type AuthorizationList struct {
	Items []Authorization
}

func (l *AuthorizationList) Add(item Authorization) {
	l.Items = append(l.Items, item)
}

func (l *AuthorizationList) Pop() *Authorization {
	if len(l.Items) == 0 {
		return nil
	}

	a := l.Items[0]
	l.Items = l.Items[1:]
	return &a
}

func (l *AuthorizationList) IsEmpty() bool {
	return len(l.Items) == 0
}
