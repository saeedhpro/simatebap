package payment

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/helper"
)

type PaymentInterface interface {
}

type PaymentListStruct struct {
	Payments     []PaymentStruct `json:"payments"`
	TotalPayment float64         `json:"total_payment"`
}

type PaymentStruct struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id"`
	UserFName   string        `json:"user_fname"`
	UserLName   string        `json:"user_lname"`
	UserKnownAs string        `json:"user_known_as,omitempty"`
	Income      bool          `json:"income"`
	Amount      float64       `json:"amount"`
	PayType     float64       `json:"paytype"`
	CheckNum    string        `json:"check_num,omitempty"`
	CheckBank   string        `json:"check_bank,omitempty"`
	CheckDate   *sql.NullTime `json:"check_date,omitempty"`
	CheckStatus bool          `json:"check_status"`
	Created     *sql.NullTime `json:"created"`
	PaidFor     int           `json:"paid_for"`
	TraceCode   string        `json:"trace_code,omitempty"`
}

type CreatePaymentStruct struct {
	UserID    int64           `json:"user_id"`
	Income    bool            `json:"income"`
	Amount    float64         `json:"amount"`
	PayType   int             `json:"paytype"`
	PaidFor   int             `json:"paid_for"`
	TraceCode string          `json:"trace_code,omitempty"`
	CheckNum  string          `json:"check_num,omitempty"`
	CheckBank string          `json:"check_bank,omitempty"`
	CheckDate helper.Datetime `json:"check_date,omitempty"`
	Info      string          `json:"info,omitempty"`
	PaidTo    string          `json:"paid_to,omitempty"`
}

type UpdatePaymentStruct struct {
	UserID    int64           `json:"user_id"`
	Income    bool            `json:"income"`
	Amount    float64         `json:"amount"`
	PayType   int             `json:"paytype"`
	PaidFor   int             `json:"paid_for"`
	TraceCode string          `json:"trace_code,omitempty"`
	CheckNum  string          `json:"check_num,omitempty"`
	CheckBank string          `json:"check_bank,omitempty"`
	CheckDate helper.Datetime `json:"check_date,omitempty"`
	Info      string          `json:"info,omitempty"`
	PaidTo    string          `json:"paid_to,omitempty"`
}
