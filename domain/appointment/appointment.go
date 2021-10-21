package appointment

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"time"
)

type CreateAppointmentRequest struct {
	UserID   int64     `json:"user_id"`
	Info     string    `json:"info"`
	StartAt  time.Time `json:"start_at"`
	CaseType string    `json:"case_type"`
	Income   float64   `json:"income"`
	IsVip    int       `json:"is_vip"`
}

type UpdateAppointmentRequest struct {
	Info               string    `json:"info"`
	StartAt            time.Time `json:"start_at"`
	CaseType           string    `json:"case_type"`
	Income             float64   `json:"income"`
	IsVip              int       `json:"is_vip"`
	FuturePrescription string    `json:"future_prescription"`
	Prescription       string    `json:"prescription"`
}

type ChangeAppointmentStatusRequest struct {
	Status int `json:"status"`
}

type SimpleAppointmentInfo struct {
	ID         int64        `json:"id"`
	CaseType   string       `json:"case_type"`
	Duration   string       `json:"duration"`
	IsVip      int64        `json:"is_vip"`
	StartAt    time.Time    `json:"start_at"`
	Status     int          `json:"status"`
	UserID     int64        `json:"user_id"`
	UserFName  string       `json:"user_fname"`
	UserLName  string       `json:"user_lname"`
	UserGender string       `json:"user_gender"`
	Mobile     string       `json:"mobile"`
	Info       string       `json:"info"`
	Income     float64      `json:"income"`
	FileID     string       `json:"file_id"`
	Price      float64      `json:"price"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
}

type UserAppointmentInfo struct {
	ID                 int64         `json:"id"`
	UserID             int64         `json:"user_id"`
	CreatedAt          *sql.NullTime `json:"created_at"`
	Info               string        `json:"info,omitempty"`
	StaffID            int64         `json:"staff_id"`
	StartAt            *sql.NullTime `json:"start_at"`
	EndAt              *sql.NullTime `json:"end_at"`
	Status             int           `json:"status"`
	DirectorID         int64         `json:"director_id"`
	UpdatedAt          *sql.NullTime `json:"updated_at"`
	Income             float64       `json:"income"`
	Subject            string        `json:"subject,omitempty"`
	CaseType           string        `json:"case_type,omitempty"`
	LaboratoryCases    string        `json:"laboratory_cases,omitempty"`
	PhotographyCases   string        `json:"photography_cases,omitempty"`
	RadiologyCases     string        `json:"radiology_cases,omitempty"`
	Prescription       string        `json:"prescription,omitempty"`
	FuturePrescription string        `json:"future_prescription,omitempty"`
	LaboratoryMsg      string        `json:"laboratory_msg,omitempty"`
	PhotographyMsg     string        `json:"photography_msg,omitempty"`
	RadiologyMsg       string        `json:"radiology_msg,omitempty"`
	OrganizationID     int64         `json:"organization_id"`
	LaboratoryID       int64         `json:"laboratory_id"`
	PhotographyID      int64         `json:"photography_id"`
	RadiologyID        int64         `json:"radiology_id"`
	LAdmissionAt       *sql.NullTime `json:"l_admission_at"`
	RAdmissionAt       *sql.NullTime `json:"r_admission_at"`
	PAdmissionAt       *sql.NullTime `json:"p_admission_at"`
	LResultAt          *sql.NullTime `json:"l_result_at"`
	RResultAt          *sql.NullTime `json:"r_result_at"`
	PResultAt          *sql.NullTime `json:"p_result_at"`
	LRndImg            string        `json:"l_rnd_img,omitempty"`
	RRndImg            string        `json:"r_rnd_img,omitempty"`
	PRndImg            string        `json:"p_rnd_img,omitempty"`
	LImgs              int           `json:"l_imgs"`
	RImgs              int           `json:"r_imgs"`
	PImgs              int           `json:"p_imgs"`
	Code               string        `json:"code,omitempty"`
	IsVip              bool          `json:"is_vip"`
	VipIntroducer      int64         `json:"vip_introducer"`
	Absence            int           `json:"absence"`
	FileID             string        `json:"file_id,omitempty"`
	Price              float64       `json:"price"`
}

type OperationInfo struct {
	ID       int64                          `json:"id"`
	UserID   int64                          `json:"user_id"`
	User     *organization.OrganizationUser `json:"user"`
	StartAt  sql.NullTime                   `json:"start_at"`
	Info     string                         `json:"info"`
	Income   float64                        `json:"income"`
	CaseType string                         `json:"case_type"`
}

type QueDetail struct {
	DefaultDuration int                     `json:"default_duration"`
	Limits          []Limit                 `json:"limits"`
	Ques            []SimpleAppointmentInfo `json:"ques"`
	Totals          []TotalLimit            `json:"totals"`
	WorkHours       *WorkHour               `json:"work_hours"`
}

type Limit struct {
	ID         int64        `json:"id"`
	Name       string       `json:"name"`
	Limitation int          `json:"limitation"`
	Total      []TotalLimit `json:"total"`
}

type TotalLimit struct {
	Total int    `json:"total"`
	Date  string `json:"date"`
}

type WorkHour struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type AcceptAppointmentRequest struct {
	FuturePrescription string   `json:"future_prescription"`
	Prescription       string   `json:"prescription"`
	PhotographyCases   []string `json:"photography_cases"`
	RadiologyCases     []string `json:"radiology_cases"`
	PhotographyID      int64    `json:"photography_id"`
	RadiologyID        int64    `json:"radiology_id"`
	PhotographyMsg     int64    `json:"photography_msg"`
	RadiologyMsg       int64    `json:"radiology_msg"`
}
