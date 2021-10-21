package pagination

import (
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
)

type OrganizationPaginationInfo struct {
	Data        []organization.OrganizationInfo `json:"data"`
	NextPage    int                             `json:"next_page"`
	PrevPage    int                             `json:"prev_page"`
	Page        int                             `json:"page"`
	HasNextPage bool                            `json:"has_next_page"`
}

type AppointmentPaginationInfo struct {
	Data        []appointment.SimpleAppointmentInfo `json:"data"`
	NextPage    int                                 `json:"next_page"`
	PrevPage    int                                 `json:"prev_page"`
	Page        int                                 `json:"page"`
	HasNextPage bool                                `json:"has_next_page"`
	PageCount   bool                                `json:"page_count"`
}
