package auth

type ResponseAccessToken struct {
	AccessToken   string        `json:"access_token"`
	ExpiresIn     int64         `json:"expires_in"`
	UserLoginInfo UserLoginInfo `json:"user"`
}

type UserLoginRequest struct {
	Tel      string `json:"tel,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserLoginInfo struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"fname"`
	LastName        string `json:"lname"`
	Tel             string `json:"tel"`
	UserGroupID     int64  `json:"user_group_id"`
	OrganizationID int64  `json:"organization_id"`
}
