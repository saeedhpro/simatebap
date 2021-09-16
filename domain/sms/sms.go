package sms

import (
	"database/sql"
	"fmt"
	"github.com/kavenegar/kavenegar-go"
)

const ApiKey = ""
const Sender = ""

type SMS struct {
	ID             int64        `json:"id"`
	UserID         int64        `json:"user_id"`
	UserFname      string       `json:"user_fname"`
	UserLname      string       `json:"user_lname"`
	StaffID        int64        `json:"staff_id"`
	StaffFname     string       `json:"staff_fname"`
	StaffLname     string       `json:"staff_lname"`
	Number         string       `json:"number"`
	Msg            string       `json:"msg"`
	Sent           bool         `json:"sent"`
	Created        sql.NullTime `json:"created"`
	Incoming       bool         `json:"incoming"`
	OrganizationID bool         `json:"organization_id"`
}

type SendSMSRequest struct {
	UserID         int64  `json:"user_id"`
	Number         string `json:"number"`
	Msg            string `json:"msg"`
	OrganizationID int64  `json:"organization_id"`
}

type DeleteSMSRequest struct {
	IDs         []int64  `json:"ids"`
}

func (s SendSMSRequest) SendSMS() {
	api := kavenegar.New(ApiKey)
	sender := Sender
	var receptor []string
	receptor = append(receptor, s.Number)
	message := s.Msg
	if res, err := api.Message.Send(sender, receptor, message, nil); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Println(err.Error())
		case *kavenegar.HTTPError:
			fmt.Println(err.Error())
		default:
			fmt.Println(err.Error())
		}
	} else {
		for _, r := range res {
			fmt.Println("MessageID 	= ", r.MessageID)
			fmt.Println("Status    	= ", r.Status)
		}
	}
}