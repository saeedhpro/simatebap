package organization

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/helper"
)

type CreateUserRequest struct {
	FirstName      string         `json:"fname" binding:"required"`
	LastName       string         `json:"lname" binding:"required"`
	Info           string         `json:"info"`
	Relation       string         `json:"relation"`
	Description    string         `json:"description"`
	FileID         string         `json:"file_id"`
	Email          sql.NullString `json:"email"`
	Gender         string         `json:"gender"`
	UserGroupId    int64          `json:"user_group_id" binding:"required"`
	OrganizationId int64          `json:"organization_id" binding:"required"`
	Tel            string         `json:"tel" binding:"required"`
	Logo           string         `json:"logo"`
	Tel1           string         `json:"tel1"`
	Nid            string         `json:"nid"`
	BirthDate      sql.NullTime   `json:"birth_date"`
	Address        string         `json:"address"`
	Introducer     string         `json:"introducer"`
	Password       string         `json:"password"`
}

type UpdateUserRequest struct {
	FirstName   string `json:"fname"`
	LastName    string `json:"lname"`
	Info        string `json:"info"`
	Description string `json:"description"`
	FileID      string `json:"file_id"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Tel         string `json:"tel"`
	Logo        string `json:"logo"`
	Tel1        string `json:"tel1"`
	Nid         string `json:"nid"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	Relation    string `json:"relation"`
}

type ChangeUserPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type OrganizationUser struct {
	ID               int64           `json:"id"`
	FirstName        string          `json:"fname"`
	LastName         string          `json:"lname"`
	Tel              string          `json:"tel"`
	UserGroupID      int64           `json:"user_group_id,omitempty"`
	Created          helper.NullTime `json:"created,omitempty"`
	LastLogin        helper.NullTime `json:"last_login,omitempty"`
	BirthDate        helper.NullTime `json:"birth_date,omitempty"`
	OrganizationID   string          `json:"organization_id,omitempty"`
	OrganizationName string          `json:"organization_name,omitempty"`
	UserGroupName    string          `json:"user_group_name,omitempty"`
	Relation         string          `json:"relation,omitempty"`
	Description      string          `json:"description,omitempty"`
}

type SimpleUserInfo struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"fname,omitempty"`
	LastName     string `json:"lname,omitempty"`
	Organization string `json:"organization,omitempty"`
}

type UserGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type LastLoginUserInfo struct {
	ID          int64           `json:"id"`
	FirstName   string          `json:"fname,omitempty"`
	LastName    string          `json:"lname,omitempty"`
	Tel         string          `json:"tel"`
	UserGroupID int64           `json:"user_group_id,omitempty"`
	LastLogin   helper.NullTime `json:"last_login,omitempty"`
}
