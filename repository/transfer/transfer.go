package transfer

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"log"
)

type Transfer struct {
	ID             int64                                `json:"id"`
	AppointmentID  int64                                `json:"appointment_id"`
	Appointment    *appointment.SimpleAppointmentInfo    `json:"appointment"`
	OrganizationID int64                                `json:"organization_id"`
	Organization   *organization.SimpleOrganizationInfo `json:"organization"`
	ToID           int64                                `json:"to_id"`
	To             *organization.SimpleUserInfo         `json:"to"`
	StaffID        int64                                `json:"staff_id"`
	Staff          *organization.SimpleUserInfo         `json:"staff"`
	Amount         float64                              `json:"amount"`
	Status         int                                  `json:"status"`
	CreatedAt      sql.NullTime                         `json:"created_at"`
}

type TransferPagination struct {
	Data       []Transfer `json:"data"`
	PagesCount int64      `json:"pages_count"`
	Page       int        `json:"page"`
}

func GetUserTransfers(id int64, page int64) *TransferPagination {
	query := "SELECT id, organization_id, to_id, staff_id, amount, created_at, status, appointment_id FROM `transfer` WHERE `to_id` = ? LIMIT 10 OFFSET ? "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	paginated := TransferPagination{}
	if err != nil {
		log.Println(err.Error())
		return &paginated
	}
	rows, err := stmt.Query(id, page - 1)
	if err != nil {
		log.Println(err.Error())
		return &paginated
	}
	transferList := []Transfer{}
	transfer := Transfer{}
	for rows.Next() {
		err = rows.Scan(
			&transfer.ID,
			&transfer.OrganizationID,
			&transfer.ToID,
			&transfer.StaffID,
			&transfer.Amount,
			&transfer.CreatedAt,
			&transfer.Status,
			&transfer.AppointmentID,
		)
		if err != nil {
			log.Println(err.Error())
			return &paginated
		}
		transfer.Organization, _ = GetOrganizationByID(transfer.OrganizationID)
		transfer.To = GetUserByID(transfer.ToID)
		transfer.Staff = GetUserByID(transfer.StaffID)
		transfer.Appointment, _ = GetAppointmentByID(transfer.AppointmentID)
		transferList = append(transferList, transfer)
	}
	count := GetTransferByIDCount(id)
	paginated.Data = transferList
	paginated.PagesCount = count
	return &paginated
}

func GetOrganizationTransfers(id int64, page int64) *TransferPagination {
	query := "SELECT id, organization_id, to_id, staff_id, amount, created_at, status, appointment_id FROM `transfer` WHERE `organization_id` = ? LIMIT 10 OFFSET ? "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	paginated := TransferPagination{}
	if err != nil {
		log.Println(err.Error())
		return &paginated
	}
	rows, err := stmt.Query(id, page - 1)
	if err != nil {
		log.Println(err.Error())
		return &paginated
	}
	transferList := []Transfer{}
	transfer := Transfer{}
	for rows.Next() {
		err = rows.Scan(
			&transfer.ID,
			&transfer.OrganizationID,
			&transfer.ToID,
			&transfer.StaffID,
			&transfer.Amount,
			&transfer.CreatedAt,
			&transfer.Status,
			&transfer.AppointmentID,
		)
		if err != nil {
			log.Println(err.Error())
			return &paginated
		}
		transfer.Organization, _ = GetOrganizationByID(transfer.OrganizationID)
		transfer.To = GetUserByID(transfer.ToID)
		transfer.Staff = GetUserByID(transfer.StaffID)
		transfer.Appointment, _ = GetAppointmentByID(transfer.AppointmentID)
		transferList = append(transferList, transfer)
	}
	count := GetTransferByIDCount(id)
	paginated.Data = transferList
	paginated.PagesCount = count
	return &paginated
}

func GetTransferByIDCount(id int64) int64 {
	query := "SELECT COUNT(*) count FROM transfer WHERE to_id = ? "
	var count int64 = 0
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count
	}
	res := stmt.QueryRow(id)
	if err = res.Err(); err != nil {
		log.Println(err.Error())
		return count
	}
	err = res.Scan(
		&count,
	)
	if err != nil {
		log.Println(err.Error())
		return count
	}
	return count
}

func GetTransferByID(id int64) (*Transfer, error) {
	query := "SELECT id, organization_id, to_id, staff_id, amount, created_at, status, appointment_id FROM transfer WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res := stmt.QueryRow(id)
	if err = res.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transfer := Transfer{}
	err = res.Scan(
		&transfer.ID,
		&transfer.OrganizationID,
		&transfer.ToID,
		&transfer.StaffID,
		&transfer.Amount,
		&transfer.CreatedAt,
		&transfer.Status,
		&transfer.AppointmentID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transfer.Organization, _ = GetOrganizationByID(transfer.OrganizationID)
	transfer.To = GetUserByID(transfer.ToID)
	transfer.Staff = GetUserByID(transfer.StaffID)
	transfer.Appointment, _ = GetAppointmentByID(transfer.AppointmentID)
	return &transfer, nil
}

func GetOrganizationByID(id int64) (*organization.SimpleOrganizationInfo, error) {
	query := "SELECT `organization`.`id` id, ifnull(`organization`.`name`, '') name FROM `organization` WHERE `organization`.`id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res := stmt.QueryRow(id)
	if err = res.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	org := organization.SimpleOrganizationInfo{}
	err = res.Scan(
		&org.ID,
		&org.Name,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &org, nil
}

func GetAppointmentByID(id int64) (*appointment.SimpleAppointmentInfo, error) {
	query := "SELECT appointment.id id, user.fname user_fname, user.lname user_lname, user.id user_id FROM appointment LEFT JOIN user ON appointment.user_id = user.id WHERE appointment.id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res := stmt.QueryRow(id)
	if err = res.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	app := appointment.SimpleAppointmentInfo{}
	err = res.Scan(
		&app.ID,
		&app.UserFName,
		&app.UserLName,
		&app.UserID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &app, nil
}

func GetUserByID(id int64) *organization.SimpleUserInfo {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetSimpleStaffQuery)
	var userInfo organization.SimpleUserInfo
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&userInfo.ID,
		&userInfo.FirstName,
		&userInfo.LastName,
		&userInfo.Organization,
	)
	if err != nil {
		return nil
	}
	return &userInfo
}

func CreateTransfer(toID int64, staffID int64, organizationID int64, appointmentID int64, amount float64) (*Transfer, error) {
	query := "INSERT INTO `transfer`(`organization_id`, `appointment_id`, `to_id`, `staff_id`, `amount`) VALUES (?,?,?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(
		organizationID,
		appointmentID,
		toID,
		staffID,
		amount,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	transfer, err := GetTransferByID(id)
	if err != nil {
		return nil, err
	}
	return transfer, nil
}