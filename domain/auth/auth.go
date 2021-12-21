package auth

import "gitlab.com/simateb-project/simateb-backend/domain/organization"

type ResponseAccessToken struct {
	AccessToken   string        `json:"access_token"`
	ExpiresIn     int64         `json:"expires_in"`
	UserLoginInfo UserLoginInfo `json:"user"`
}

type UserLoginRequest struct {
	Tel      string `json:"tel" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginInfo struct {
	ID             int64                              `json:"id"`
	FirstName      string                             `json:"fname"`
	LastName       string                             `json:"lname"`
	Tel            string                             `json:"tel"`
	UserGroupID    int64                              `json:"user_group_id"`
	OrganizationID int64                              `json:"organization_id"`
	Profession     *organization.SimpleProfessionInfo `json:"profession"`
}
