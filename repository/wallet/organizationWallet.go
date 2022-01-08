package wallet

import (
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
	"strconv"
)

type WalletHistoriesPaginationInfo struct {
	Data        []WalletHistoryStruct `json:"data"`
	NextPage    int                          `json:"next_page"`
	PrevPage    int                          `json:"prev_page"`
	Page        int                          `json:"page"`
	HasNextPage bool                         `json:"has_next_page"`
	PagesCount  int                          `json:"pages_count"`
}

func GetOrganizationAllWalletHistories(userId int64, start_date string, end_date string, q string, page string) (WalletHistoriesPaginationInfo, error) {
	histories := []WalletHistoryStruct{}
	var history WalletHistoryStruct
	paginated := WalletHistoriesPaginationInfo{
		Data: histories,
	}
	query := "SELECT wallet_histories.`id` id, ifnull(wallet_histories.`user_id`, 0) user_id, wallet_histories.`owner_id` owner_id, wallet_histories.`balance` balance, wallet_histories.`created_at` created_at, ifnull(wallet_histories.`updated_at`, '') updated_at, wallet_histories.`status` status, wallet_histories.`type` type, ifnull(wallet_histories.`sheba`, '') sheba, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname from wallet_histories LEFT JOIN user ON wallet_histories.owner_id = user.id where type = 'withdraw' AND (wallet_histories.user_id = ? or wallet_histories.owner_id = ?) "
	var values []interface{}
	values = append(values, userId)
	values = append(values, userId)
	if start_date != "" && start_date != "null" && start_date != "undefined" {
		query += " AND wallet_histories.created_at >= ? "
		values = append(values, start_date)
	}
	if end_date != "" && end_date != "null" && end_date != "undefined" {
		query += " AND wallet_histories.created_at <= ? "
		values = append(values, end_date)
	}
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	if page != "" && page != "null" && page != "undefined" {
		query += " LIMIT 10 OFFSET ? "
		values = append(values, page)
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return paginated, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println(err.Error())
		return paginated, err
	}
	for rows.Next() {
		err = rows.Scan(
			&history.ID,
			&history.UserID,
			&history.OwnerID,
			&history.Balance,
			&history.CreatedAt,
			&history.UpdatedAt,
			&history.Status,
			&history.Type,
			&history.Sheba,
			&history.FName,
			&history.LName,
		)
		if err != nil {
			log.Println(err.Error())
			return paginated, err
		}
		histories = append(histories, history)
	}
	p, err := strconv.Atoi(page)
	count := 0
	count, _ = GetOrganizationAllWalletHistoriesCount(userId, start_date, end_date, q)
	paginated = WalletHistoriesPaginationInfo{
		Data:       histories,
		Page:       p,
		PagesCount: count,
	}
	return paginated, nil
}

func GetOrganizationAllWalletHistoriesCount(userId int64, start_date string, end_date string, q string) (int, error) {
	query := "SELECT COUNT(*) from wallet_histories LEFT JOIN user ON wallet_histories.owner_id = user.id where 'withdraw' AND (wallet_histories.user_id = ? or wallet_histories.owner_id = ?) "
	var values []interface{}
	count := 0
	values = append(values, userId)
	values = append(values, userId)
	if start_date != "" && start_date != "null" && start_date != "undefined" {
		query += " AND wallet_histories.created_at >= ? "
		values = append(values, start_date)
	}
	if end_date != "" && end_date != "null" && end_date != "undefined" {
		query += " AND wallet_histories.created_at <= ? "
		values = append(values, end_date)
	}
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func GetOrganizationWalletHistoriesSum(userId int64, start_date string, end_data string) (int64, error) {
	var sum int64
	query := "select s.sum from (SELECT SUM(balance) as sum FROM wallet_histories WHERE owner_id = ? and type = 'withdraw' and status = 2 and created_at between ? and ?) as s where s.sum is not null"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0, err
	}
	row := stmt.QueryRow(userId, start_date, end_data)
	if row.Err() != nil {
		return 0, nil
	}
	err = row.Scan(
		&sum,
	)
	if err != nil {
		log.Println(err.Error(), " :err: ")
		return 0, err
	}
	return sum, nil
}

func GetOrganizationWalletHistoriesDays(userId int64, start_date string, end_data string) ([]WalletHistoryStruct, error) {
	query := "SELECT `id`, `user_id`, `owner_id`, `balance`, `created_at`, `updated_at`, `status`, `type`, `sheba` FROM wallet_histories WHERE (user_id = ? or owner_id = ?) and created_at between ? and ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	histories := []WalletHistoryStruct{}
	var history WalletHistoryStruct
	if err != nil {
		log.Println(err.Error(), "prepare")
		return histories, err
	}
	rows, err := stmt.Query(userId, userId, start_date, end_data)
	if err != nil {
		log.Println(err.Error())
		return histories, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&history.ID,
			&history.UserID,
			&history.OwnerID,
			&history.Balance,
			&history.CreatedAt,
			&history.UpdatedAt,
			&history.Status,
			&history.Type,
			&history.Sheba,
		)
		if err != nil {
			log.Println(err.Error())
			return histories, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}


func GetUserAllWalletHistories(userId int64, start_date string, end_date string, q string, page string) (WalletHistoriesPaginationInfo, error) {
	histories := []WalletHistoryStruct{}
	var history WalletHistoryStruct
	paginated := WalletHistoriesPaginationInfo{
		Data: histories,
	}
	query := "SELECT wallet_histories.`id` id, ifnull(wallet_histories.`user_id`, 0) user_id, wallet_histories.`owner_id` owner_id, wallet_histories.`balance` balance, wallet_histories.`created_at` created_at, ifnull(wallet_histories.`updated_at`, '') updated_at, wallet_histories.`status` status, wallet_histories.`type` type, ifnull(wallet_histories.`sheba`, '') sheba, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname from wallet_histories LEFT JOIN user ON wallet_histories.owner_id = user.id where type = 'withdraw' AND (wallet_histories.user_id = ? or wallet_histories.owner_id = ?) "
	var values []interface{}
	values = append(values, userId)
	values = append(values, userId)
	if start_date != "" && start_date != "null" && start_date != "undefined" {
		query += " AND wallet_histories.created_at >= ? "
		values = append(values, start_date)
	}
	if end_date != "" && end_date != "null" && end_date != "undefined" {
		query += " AND wallet_histories.created_at <= ? "
		values = append(values, end_date)
	}
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	if page != "" && page != "null" && page != "undefined" {
		query += " LIMIT 10 OFFSET ? "
		values = append(values, page)
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return paginated, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println(err.Error())
		return paginated, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&history.ID,
			&history.UserID,
			&history.OwnerID,
			&history.Balance,
			&history.CreatedAt,
			&history.UpdatedAt,
			&history.Status,
			&history.Type,
			&history.Sheba,
			&history.FName,
			&history.LName,
		)
		if err != nil {
			log.Println(err.Error())
			return paginated, err
		}
		histories = append(histories, history)
	}
	p, err := strconv.Atoi(page)
	count := 0
	count, _ = GetUserAllWalletHistoriesCount(userId, start_date, end_date, q)
	paginated = WalletHistoriesPaginationInfo{
		Data:       histories,
		Page:       p,
		PagesCount: count,
	}
	return paginated, nil
}
func GetUserAllWalletHistoriesCount(userId int64, start_date string, end_date string, q string) (int, error) {
	query := "SELECT COUNT(*) from wallet_histories LEFT JOIN user ON wallet_histories.owner_id = user.id where 'withdraw' AND (wallet_histories.user_id = ? or wallet_histories.owner_id = ?) "
	var values []interface{}
	count := 0
	values = append(values, userId)
	values = append(values, userId)
	if start_date != "" && start_date != "null" && start_date != "undefined" {
		query += " AND wallet_histories.created_at >= ? "
		values = append(values, start_date)
	}
	if end_date != "" && end_date != "null" && end_date != "undefined" {
		query += " AND wallet_histories.created_at <= ? "
		values = append(values, end_date)
	}
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func GetUserWalletHistoriesSum(userId int64, start_date string, end_data string) (int64, error) {
	var sum int64
	query := "select s.sum from (SELECT SUM(balance) as sum FROM wallet_histories WHERE owner_id = ? and type = 'withdraw' and status = 2 and created_at between ? and ?) as s where s.sum is not null"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0, err
	}
	row := stmt.QueryRow(userId, start_date, end_data)
	if row.Err() != nil {
		return 0, nil
	}
	err = row.Scan(
		&sum,
	)
	if err != nil {
		log.Println(err.Error(), " :err: ")
		return 0, err
	}
	return sum, nil
}

func GetUserWalletHistoriesDays(userId int64, start_date string, end_data string) ([]WalletHistoryStruct, error) {
	query := "SELECT `id`, `user_id`, `owner_id`, `balance`, `created_at`, `updated_at`, `status`, `type`, `sheba` FROM wallet_histories WHERE (user_id = ? or owner_id = ?) and created_at between ? and ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	histories := []WalletHistoryStruct{}
	var history WalletHistoryStruct
	if err != nil {
		log.Println(err.Error(), "prepare")
		return histories, err
	}
	rows, err := stmt.Query(userId, userId, start_date, end_data)
	if err != nil {
		log.Println(err.Error())
		return histories, err
	}
	for rows.Next() {
		err = rows.Scan(
			&history.ID,
			&history.UserID,
			&history.OwnerID,
			&history.Balance,
			&history.CreatedAt,
			&history.UpdatedAt,
			&history.Status,
			&history.Type,
			&history.Sheba,
		)
		if err != nil {
			log.Println(err.Error())
			return histories, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}
