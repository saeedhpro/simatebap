package wallet

type WalletHistoryStruct struct {
	ID        int64   `json:"id"`
	UserID    *int64  `json:"user_id"`
	OwnerID   int64   `json:"owner_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Status    int64   `json:"status"`
	Type      string  `json:"type"`
	Sheba     string  `json:"sheba"`
	FName     string  `json:"fname"`
	LName     string  `json:"lname"`
}
