package wallet

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type WalletInterface interface {
	User() (*WalletUserStruct, error)
	Increase(amount float64) bool
	Decrease(amount float64) bool
	Create() bool
}

type ChangeUserWalletBalance struct {
	Amount float64 `json:"amount"`
}

type WalletUserStruct struct {
	ID        string `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
}

type WalletStruct struct {
	ID      int64         `json:"id"`
	UserID  int64         `json:"user_id"`
	Balance float64       `json:"balance"`
	Created *sql.NullTime `json:"created"`
}

func (w *WalletStruct) User() (*WalletUserStruct, error) {
	var user WalletUserStruct
	query := "SElECT id, fname, lname FROM user WHERE user.id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}
	result := stmt.QueryRow(w.UserID)
	err = result.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (w *WalletStruct) Increase(amount float64) (bool, float64) {
	query := "UPDATE `wallet` SET `balance`= ?  WHERE user_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false, 0
	}
	newBalance := w.Balance + amount
	_, err = stmt.Exec(newBalance, w.UserID)
	if err != nil {
		return false, 0
	}
	return true, newBalance
}

func (w *WalletStruct) Decrease(amount float64, force bool) (bool, float64) {
	query := "UPDATE `wallet` SET `balance`= ?  WHERE user_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false, 0
	}
	if !force && w.Balance < amount {
		return false, 0
	}
	newBalance := w.Balance - amount
	_, err = stmt.Exec(newBalance, w.UserID)
	if err != nil {
		return false, 0
	}
	return true, newBalance
}

func (w *WalletStruct) SetBalance(amount float64) bool {
	query := "UPDATE `wallet` SET `balance`= ?  WHERE user_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(amount, w.UserID)
	if err != nil {
		return false
	}
	return true
}

func (w *WalletStruct) Create(t string) bool {
	query := "INSERT INTO `wallet` (`user_id`, `balance`, `type`) VALUES (?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(w.UserID, 0, t)
	if err != nil {
		return false
	}
	return true
}

func GetWallet(userID int64, t string) *WalletStruct {
	wallet := WalletStruct{
		UserID: userID,
	}
	var count int
	user, err := wallet.User()
	if err != nil {
		log.Println(err.Error())
	}
	if user != nil {
		query := "SELECT COUNT(*) count FROM wallet WHERE user_id = ? AND type = ?"
		stmt, err := repository.DBS.MysqlDb.Prepare(query)
		if err != nil {
			return nil
		}
		result := stmt.QueryRow(userID, t)
		err = result.Scan(
			&count,
		)
		if err != nil {
			return nil
		}
		if count == 0 {
			wallet.Create(t)
		}
	}
	query := "SElECT id, user_id, balance, created FROM wallet WHERE user_id = ? AND type = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(userID, t)
	err = result.Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Created,
	)
	if err != nil {
		return nil
	}
	return &wallet
}
