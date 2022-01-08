package sms

import (
	"gitlab.com/simateb-project/simateb-backend/domain/sms"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"log"
)

func SendSMS(sendSMSRequest sms.SendSMSRequest, staffID int64) (*string, error) {
	var sent bool
	sent, r, error := sendSMSRequest.SendSMS()
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateSMSQuery)
	if err != nil {
		log.Println(err.Error())
		return r, err
	}
	result, err := stmt.Exec(
		&sendSMSRequest.UserID,
		staffID,
		&sendSMSRequest.Number,
		&sendSMSRequest.Msg,
		sent,
		false,
		&sendSMSRequest.OrganizationID,
	)
	if err != nil {
		log.Println(err.Error())
		return r, err
	}
	res, err := result.LastInsertId()
	log.Println(res)
	if err != nil {
		log.Println(err.Error())
		return r, err
	}
	if error != nil {
		return nil, error
	}
	return nil, nil
}

func GetList(orgID int64) ([]sms.SMS, error) {
	var query = "SELECT "
	var smsList []sms.SMS
	var sms sms.SMS
	log.Println(query)
	log.Println(sms)
	return smsList, nil
}

func GetListForAdmin() ([]sms.SMS, error) {
	var smsList []sms.SMS

	return smsList, nil
}
