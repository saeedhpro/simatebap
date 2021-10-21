package paymentController

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	payment2 "gitlab.com/simateb-project/simateb-backend/domain/payment"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strings"
)

type PaymentControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetPaymentList(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type PaymentControllerStruct struct {
}

func NewPaymentController() PaymentControllerInterface {
	x := &PaymentControllerStruct{
	}
	return x
}

func (pc *PaymentControllerStruct) Create(c *gin.Context) {
	var request payment2.CreatePaymentStruct
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	query := "INSERT INTO `payment`(`user_id`, `income`, `amount`, `paytype`, `paid_for`, `trace_code`, `check_num`, `check_bank`, `check_date`, `info`, `paid_to`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		&request.UserID,
		&request.Income,
		&request.Amount,
		&request.PayType,
		&request.PaidFor,
		&request.TraceCode,
		&request.CheckNum,
		&request.CheckBank,
		&request.CheckDate.Time,
		&request.Info,
		&request.PaidTo,
	)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, true)
}

func (pc *PaymentControllerStruct) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	query := "SELECT payment.id id, user.id user_id, user.fname user_fname, user.lname user_name, ifnull(user.known_as, '') user_known_as, payment.income income, payment.amount amount, payment.paytype paytype, ifnull(payment.check_num, '') check_num, ifnull(payment.check_bank, '') check_bank, ifnull(payment.check_date, '') check_date, payment.check_status check_status, payment.created created FROM payment LEFT JOIN user ON payment.user_id = user.id WHERE payment.id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var payment payment2.PaymentStruct
	result := stmt.QueryRow(id)
	err = result.Scan(
		&payment.ID,
		&payment.UserFName,
		&payment.UserFName,
		&payment.UserLName,
		&payment.UserKnownAs,
		&payment.Income,
		&payment.Amount,
		&payment.PayType,
		&payment.CheckNum,
		&payment.CheckBank,
		&payment.CheckDate,
		&payment.CheckStatus,
		&payment.Created,
	)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "یافت نشد")
			return
		}
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (pc *PaymentControllerStruct) GetPaymentList(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	query := "SELECT payment.id id, user.id user_id, user.fname user_fname, user.lname user_name, ifnull(user.known_as, '') user_known_as, payment.income income, payment.amount amount, payment.paytype paytype, ifnull(payment.check_num, '') check_num, ifnull(payment.check_bank, '') check_bank, payment.check_date check_date, payment.check_status check_status, payment.created created, payment.paid_for paid_for, payment.trace_code trace_code FROM payment LEFT JOIN user ON payment.user_id = user.id WHERE payment.user_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	rows, error := stmt.Query(id)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	var payments []payment2.PaymentStruct
	var payment payment2.PaymentStruct
	for rows.Next() {
		err = rows.Scan(
			&payment.ID,
			&payment.UserFName,
			&payment.UserFName,
			&payment.UserLName,
			&payment.UserKnownAs,
			&payment.Income,
			&payment.Amount,
			&payment.PayType,
			&payment.CheckNum,
			&payment.CheckBank,
			&payment.CheckDate,
			&payment.CheckStatus,
			&payment.Created,
			&payment.PaidFor,
			&payment.TraceCode,
		)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data": err.Error(),
			})
			return
		}
		payments = append(payments, payment)
	}
	var total float64
	query = "SELECT SUM(payment.amount) total FROM payment LEFT JOIN user ON payment.user_id = user.id WHERE payment.user_id = ?"
	stmt, err = repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&total,
	)
	if err != nil {
		log.Println(err.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	paymentList := payment2.PaymentListStruct{
		Payments:     payments,
		TotalPayment: total,
	}
	c.JSON(http.StatusOK, paymentList)
}

func (pc *PaymentControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `payment` SET"
	var values []interface{}
	var columns []string
	userId := c.Param("id")
	if userId == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var request payment2.UpdatePaymentStruct
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getAppointmentUpdateColumns(&request, &columns, &values)
	columnsString := strings.Join(columns, ",")
	updateUserQuery += columnsString
	updateUserQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateUserQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, userId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func getAppointmentUpdateColumns(o *payment2.UpdatePaymentStruct, columns *[]string, values *[]interface{}) {
	*columns = append(*columns, " `income` = ? ")
	*values = append(*values, o.Income)
	*columns = append(*columns, " `amount` = ? ")
	*values = append(*values, o.Amount)
	*columns = append(*columns, " `paytype` = ? ")
	*values = append(*values, o.PayType)
	*columns = append(*columns, " `paid_for` = ? ")
	*values = append(*values, o.PaidFor)
	*columns = append(*columns, " `paid_to` = ? ")
	*values = append(*values, o.PaidTo)
	*columns = append(*columns, " `trace_code` = ? ")
	*values = append(*values, o.TraceCode)
	if o.CheckNum != "" {
		*columns = append(*columns, " `check_num` = ? ")
		*values = append(*values, o.CheckNum)
	}
	if o.CheckBank != "" {
		*columns = append(*columns, " `check_bank` = ? ")
		*values = append(*values, o.CheckBank)
	}
	if o.CheckDate != (helper.Datetime{}) {
		*columns = append(*columns, " `check_date` = ? ")
		*values = append(*values, o.CheckDate.Time)
	}
	if o.Info != "" {
		*columns = append(*columns, " `info` = ? ")
		*values = append(*values, o.Info)
	}
}

func (pc *PaymentControllerStruct) Delete(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	query := "DELETE FROM `payment` WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	stmt.QueryRow(userID)
	c.JSON(200, nil)
}
