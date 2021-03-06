package sms

import (
	"database/sql"
	"github.com/kavenegar/kavenegar-go"
	"log"
)

const ApiKey = "68594136704C39444A5474653364346B387A64534D5132664F6F36646F4E3875"

//const Sender = "0018018949161"
const Sender = "10008663"

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

type SendNotificationRequest struct {
	UserID       int64  `json:"user_id"`
	StaffID      int64  `json:"staff_id"`
	ActionUrl    string `json:"action_url"`
	CloseOnClick int    `json:"close_on_click"`
	Content      string `json:"content"`
	Title        string `json:"title"`
	IDs          string `json:"ids"`
	Type         string `json:"type"`
}

type DeleteSMSRequest struct {
	IDs []int64 `json:"ids"`
}

func (s SendSMSRequest) SendSMS() (bool, *string, error) {
	sender := Sender
	var receptor []string
	receptor = append(receptor, s.Number)
	message := s.Msg
	send, res, error := SendByPackage(sender, receptor, message)
	return send, res, error
}

func SendByPackage(sender string, receptor []string, message string) (bool, *string, error) {
	api := kavenegar.New(ApiKey)
	if res, err := api.Message.Send(sender, receptor, message, nil); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			log.Println(err.Error())
			break
		case *kavenegar.HTTPError:
			log.Println(err.Error())
			break
		default:
			log.Println(err.Error())
			break
		}
		r := err.Error()
		return false, &r, err
	} else {
		for _, r := range res {
			log.Println("MessageID 	= ", r.MessageID)
			log.Println("Status    	= ", r.Status)
		}
		return true, nil, nil
	}
}
