package domain

type Response struct {
	Accounts []Account `json:"accounts,omitempty"`
	Message  string    `json:"msg,omitempty"`
	Error    string    `json:"error,omitempty"`
}

type Account struct {
	ID     string `json:"id,omitempty"`
	Ledger uint32 `json:"ledger,omitempty"`
	Code   uint16 `json:"code,omitempty"`
}

type TransferAmount struct {
	SenderAccountID    string  `json:"sender_account_id,omitempty"`
	RecipientAccountID string  `json:"recipient_account_id,omitempty"`
	Amount             float32 `json:"amount,omitempty"`
}
