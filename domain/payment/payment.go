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
	UserKnownAs string        `json:"user_known_as"`
	Income      bool          `json:"income"`
	Amount      float64       `json:"amount"`
	PayType     float64       `json:"paytype"`
	CheckNum    string        `json:"check_num"`
	CheckBank   string        `json:"check_bank"`
	CheckDate   *sql.NullTime `json:"check_date"`
	CheckStatus bool          `json:"check_status"`
	Created     *sql.NullTime `json:"created"`
	PaidFor     int           `json:"paid_for"`
	TraceCode   string        `json:"trace_code"`
}

type CreatePaymentStruct struct {
	UserID    int64   `json:"user_id"`
	Income    bool    `json:"income"`
	Amount    float64 `json:"amount"`
	PayType   int     `json:"paytype"`
	PaidFor   int     `json:"paid_for"`
	TraceCode string  `json:"trace_code"`
	CheckNum  string  `json:"check_num"`
	CheckBank string  `json:"check_bank"`
	CheckDate string  `json:"check_date"`
	Info      string  `json:"info"`
	PaidTo    string  `json:"paid_to"`
}

type UpdatePaymentStruct struct {
	UserID    int64           `json:"user_id"`
	Income    bool            `json:"income"`
	Amount    float64         `json:"amount"`
	PayType   int             `json:"paytype"`
	PaidFor   int             `json:"paid_for"`
	TraceCode string          `json:"trace_code"`
	CheckNum  string          `json:"check_num"`
	CheckBank string          `json:"check_bank"`
	CheckDate helper.Datetime `json:"check_date"`
	Info      string          `json:"info"`
	PaidTo    string          `json:"paid_to"`
}
