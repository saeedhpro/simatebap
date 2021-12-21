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
	BirthDate      string         `json:"birth_date"`
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
	ID               int64                 `json:"id"`
	Logo             string                `json:"logo"`
	AppCode          string                `json:"appcode"`
	FirstName        string                `json:"fname"`
	LastName         string                `json:"lname"`
	Tel              string                `json:"tel"`
	UserGroupID      int                   `json:"user_group_id"`
	Created          *sql.NullTime         `json:"created"`
	LastLogin        *sql.NullTime         `json:"last_login"`
	BirthDate        *sql.NullTime         `json:"birth_date"`
	Birth            int                   `json:"birth"`
	OrganizationID   string                `json:"organization_id"`
	OrganizationName string                `json:"organization_name"`
	UserGroupName    string                `json:"user_group_name"`
	Relation         string                `json:"relation"`
	Description      string                `json:"description"`
	Info             string                `json:"info"`
	Tel1             string                `json:"tel1"`
	Nid              string                `json:"nid"`
	Address          string                `json:"address"`
	Introducer       string                `json:"introducer"`
	Gender           string                `json:"gender"`
	FileID           string                `json:"file_id"`
	Profession       *SimpleProfessionInfo `json:"profession"`
}

type LastLoginUser struct {
	ID               int64         `json:"id"`
	UserFirstName    string        `json:"user_fname"`
	UserLastName     string        `json:"user_lname"`
	Tel              string        `json:"tel"`
	LastLogin        *sql.NullTime `json:"last_login"`
	OrganizationID   string        `json:"organization_id"`
	OrganizationName string        `json:"organization_name"`
	UserGroupName    string        `json:"user_group_name"`
}

type SimpleUserInfo struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"fname"`
	LastName     string `json:"lname"`
	Organization string `json:"organization"`
}

type UserGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type LastLoginUserInfo struct {
	ID          int64           `json:"id"`
	FirstName   string          `json:"fname"`
	LastName    string          `json:"lname"`
	Tel         string          `json:"tel"`
	UserGroupID int64           `json:"user_group_id"`
	LastLogin   helper.NullTime `json:"last_login"`
}

type OrganizationUserPaginate struct {
	Data        []OrganizationUser `json:"data"`
	NextPage    int                `json:"next_page"`
	PrevPage    int                `json:"prev_page"`
	Page        int                `json:"page"`
	HasNextPage bool               `json:"has_next_page"`
	PagesCount  int                `json:"pages_count"`
}
