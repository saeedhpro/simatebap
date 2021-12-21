package holiday

import "time"

type CreateHolidayRequest struct {
	OrganizationID int64  `json:"organization_id"`
	HDate          string `json:"hdate"`
	Title          string `json:"title"`
}

type UpdateHolidayRequest struct {
	HDate          time.Time `json:"hdate"`
	Title          string    `json:"title"`
	OrganizationID int64     `json:"organization_id"`
}

type HolidayInfo struct {
	ID               int64     `json:"id"`
	OrganizationID   int64     `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	HDate            time.Time `json:"hdate"`
	Title            string    `json:"title"`
}
