package dashboard

import "gitlab.com/simateb-project/simateb-backend/domain/organization"

type DashboardData struct {
	Count                  int                              `json:"count"`
	UnknownGenderUserCount int                              `json:"unknown_gender_user_count"`
	FemaleGenderUserCount  int                              `json:"female_gender_user_count"`
	MaleGenderUserCount    int                              `json:"male_gender_user_count"`
	LastLoginUsers         []organization.LastLoginUserInfo `json:"last_login_users"`
}
