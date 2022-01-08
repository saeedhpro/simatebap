package user

import (
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/repository/transfer"
	"log"
	"time"
)

func GetUserByID(userId int64) (*organization.OrganizationUser, error) {
	var user organization.OrganizationUser
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetUserOrganizationQuery)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return nil, err
	}
	result := stmt.QueryRow(userId)
	err = result.Scan(
		&user.ID,
		&user.AppCode,
		&user.Logo,
		&user.FirstName,
		&user.LastName,
		&user.LastLogin,
		&user.Created,
		&user.Tel,
		&user.UserGroupID,
		&user.UserGroupName,
		&user.BirthDate,
		&user.OrganizationID,
		&user.OrganizationName,
		&user.Relation,
		&user.Description,
		&user.Info,
		&user.Tel1,
		&user.Nid,
		&user.Address,
		&user.Introducer,
		&user.Gender,
		&user.FileID,
	)
	user.Profession = organizationController.GetProfession(user.OrganizationID)
	if user.BirthDate != nil {
		year, _, _, _, _, _ := helper.TimeDiff(user.BirthDate.Time, time.Now())
		user.Birth = year
	}
	if err != nil {
		log.Println(err.Error(), " :err: ")
	}
	return &user, nil
}

type Withdraw struct {
	ID        int64                        `json:"id"`
	OwnerID   int64                        `json:"owner_id"`
	UserID    int64                        `json:"user_id"`
	User      *organization.SimpleUserInfo `json:"user"`
	Owner     *organization.SimpleUserInfo `json:"owner"`
	Balance   float64                      `json:"balance"`
	Sheba     string                       `json:"sheba"`
	Status    int                          `json:"status"`
	CreatedAt string                       `json:"created_at"`
}

type WithdrawPaginated struct {
	Data       []Withdraw `json:"data"`
	PagesCount int        `json:"pages_count"`
}

type CreateWithdrawRequest struct {
	OwnerID int64   `json:"owner_id"`
	Balance float64 `json:"balance"`
	Sheba   string  `json:"sheba"`
}
type AcceptOrRejectWithdrawRequest struct {
	ID       int64 `json:"id"`
	Accepted bool  `json:"accepted"`
}

func CreateWithdraw(request CreateWithdrawRequest) (bool, error) {
	query := "INSERT INTO `wallet_histories`(`owner_id`, `balance`, `status`, `type`, `sheba`) VALUES (?,?,?,?,?)"

	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		&request.OwnerID,
		&request.Balance,
		1,
		"withdraw",
		&request.Sheba,
	)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	log.Println(res.LastInsertId())
	return true, nil
}

func AcceptOrRejectWithdraw(request AcceptOrRejectWithdrawRequest) (bool, error) {
	query := "UPDATE `wallet_histories` SET `status` = ? WHERE id = ?"

	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	defer stmt.Close()
	status := 0
	if request.Accepted {
		status = 2
	}
	_, err = stmt.Exec(
		status,
		request.ID,
	)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

func GetWalletHistoriesForAdmin(page string) (*WithdrawPaginated, error) {
	query := "SELECT `id`, ifnull(`user_id`, 0) `user_id`,ifnull(`owner_id`, 0) `owner_id`, `balance`, `created_at`, `status`, `sheba` FROM `wallet_histories` WHERE `type` = 'withdraw' LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	pageinated := WithdrawPaginated{}
	if err != nil {
		log.Println(err.Error())
		return &pageinated, err
	}
	rows, err := stmt.Query(page)
	if err != nil {
		log.Println(err.Error())
		return &pageinated, err
	}
	list := []Withdraw{}
	var w Withdraw
	for rows.Next() {
		err = rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OwnerID,
			&w.Balance,
			&w.CreatedAt,
			&w.Status,
			&w.Sheba,
		)
		if err != nil {
			log.Println(err.Error())
			return &pageinated, err
		}
		w.User = transfer.GetUserByID(w.UserID)
		w.Owner = transfer.GetUserByID(w.OwnerID)
		list = append(list, w)
	}
	pageinated.Data = list
	count := 0
	count = GetWalletHistoriesForAdminCount(page)
	pageinated.PagesCount = count
	return &pageinated, nil
}

func GetWalletHistoriesForAdminCount(page string) int {
	count := 0
	query := "SELECT COUNT(*) count FROM `wallet_histories` WHERE `type` = 'withdraw'"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count
	}
	row := stmt.QueryRow(page)
	if err = row.Err(); err != nil {
		log.Println(err.Error())
		return count
	}
	err = row.Scan(
		&count,
	)
	if err != nil {
		log.Println(err.Error())
		return count
	}
	return count
}
