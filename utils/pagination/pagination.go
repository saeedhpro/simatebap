package pagination

import (
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	sms2 "gitlab.com/simateb-project/simateb-backend/domain/sms"
	"gitlab.com/simateb-project/simateb-backend/repository/wallet"
)

type OrganizationPaginationInfo struct {
	Data        []organization.OrganizationInfo `json:"data"`
	NextPage    int                             `json:"next_page"`
	PrevPage    int                             `json:"prev_page"`
	Page        int                             `json:"page"`
	HasNextPage bool                            `json:"has_next_page"`
	PagesCount  int                             `json:"pages_count"`
}

type OrganizationUserPaginationInfo struct {
	Data        []organization.OrganizationUser `json:"data"`
	NextPage    int                             `json:"next_page"`
	PrevPage    int                             `json:"prev_page"`
	Page        int                             `json:"page"`
	HasNextPage bool                            `json:"has_next_page"`
	PagesCount  int                             `json:"pages_count"`
}

type OrganizationAppointmentPaginationInfo struct {
	Data        []appointment.UserAppointmentInfo `json:"data"`
	NextPage    int                               `json:"next_page"`
	PrevPage    int                               `json:"prev_page"`
	Page        int                               `json:"page"`
	HasNextPage bool                              `json:"has_next_page"`
	PagesCount  int                               `json:"pages_count"`
}

type OrganizationAbout struct {
	Text1  string `json:"text1"`
	Text2  string `json:"text2"`
	Text3  string `json:"text3"`
	Text4  string `json:"text4"`
	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
	Image3 string `json:"image3"`
	Image4 string `json:"image4"`
}

type OrganizationAboutRequest struct {
	Text1  string `json:"text1"`
	Text2  string `json:"text2"`
	Text3  string `json:"text3"`
	Text4  string `json:"text4"`
	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
	Image3 string `json:"image3"`
	Image4 string `json:"image4"`
}

type OrganizationWorkTimeStruct struct {
	WorkHourStart string `json:"work_hour_start"`
	WorkHourEnd   string `json:"work_hour_end"`
}

type WalletHistoriesPaginationInfo struct {
	Data        []wallet.WalletHistoryStruct `json:"data"`
	NextPage    int                          `json:"next_page"`
	PrevPage    int                          `json:"prev_page"`
	Page        int                          `json:"page"`
	HasNextPage bool                         `json:"has_next_page"`
	PagesCount  int                          `json:"pages_count"`
}
type SMSPaginationInfo struct {
	Data        []sms2.SMS `json:"data"`
	NextPage    int        `json:"next_page"`
	PrevPage    int        `json:"prev_page"`
	Page        int        `json:"page"`
	HasNextPage bool       `json:"has_next_page"`
	PagesCount  int        `json:"pages_count"`
}

type AppointmentPaginationInfo struct {
	Data            []appointment.SimpleAppointmentInfo `json:"data"`
	HasNextPage     bool                                `json:"has_next_page"`
	HasPreviousPage bool                                `json:"has_previous_page"`
	PagesCount      int                                 `json:"pages_count"`
}

type SendDocAppointmentPaginationInfo struct {
	Data            []appointment.UserAppointmentInfo `json:"data"`
	HasNextPage     bool                              `json:"has_next_page"`
	HasPreviousPage bool                              `json:"has_previous_page"`
	PagesCount      int                               `json:"pages_count"`
}
