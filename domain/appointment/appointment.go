package appointment

import (
	"database/sql"
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
	UpdatedAt  sql.NullTime `json:"updated_at"`
}

type UserAppointmentInfo struct {
	ID                 int64        `json:"id"`
	UserID             int64        `json:"user_id"`
	CreatedAt          sql.NullTime `json:"created_at"`
	Info               string       `json:"info"`
	StaffID            int64        `json:"staff_id"`
	StartAt            time.Time    `json:"start_at"`
	EndAt              sql.NullTime `json:"end_at"`
	Status             int          `json:"status"`
	DirectorID         int64        `json:"director_id"`
	UpdatedAt          sql.NullTime `json:"updated_at"`
	Income             float64      `json:"income"`
	Subject            string       `json:"subject"`
	CaseType           string       `json:"case_type"`
	LaboratoryCases    string       `json:"laboratory_cases"`
	PhotographyCases   string       `json:"photography_cases"`
	RadiologyCases     string       `json:"radiology_cases"`
	Prescription       string       `json:"prescription"`
	FuturePrescription string       `json:"future_prescription"`
	LaboratoryMsg      string       `json:"laboratory_msg"`
	PhotographyMsg     string       `json:"photography_msg"`
	RadiologyMsg       string       `json:"radiology_msg"`
	OrganizationID     int64        `json:"organization_id"`
	LaboratoryID       int64        `json:"laboratory_id"`
	PhotographyID      int64        `json:"photography_id"`
	RadiologyID        int64        `json:"radiology_id"`
	LAdmissionAt       sql.NullTime `json:"l_admission_at"`
	RAdmissionAt       sql.NullTime `json:"r_admission_at"`
	PAdmissionAt       sql.NullTime `json:"p_admission_at"`
	LResultAt          sql.NullTime `json:"l_result_at"`
	RResultAt          sql.NullTime `json:"r_result_at"`
	PResultAt          sql.NullTime `json:"p_result_at"`
	LRndImg            string       `json:"l_rnd_img"`
	RRndImg            string       `json:"r_rnd_img"`
	PRndImg            string       `json:"p_rnd_img"`
	LImgs              int          `json:"l_imgs"`
	RImgs              int          `json:"r_imgs"`
	PImgs              int          `json:"p_imgs"`
	Code               string       `json:"code"`
	IsVip              bool         `json:"is_vip"`
	VipIntroducer      int64        `json:"vip_introducer"`
	Absence            int          `json:"absence"`
}

type OperationInfo struct {
	ID       int64        `json:"id"`
	UserID   int64        `json:"user_id"`
	StartAt  sql.NullTime `json:"start_at"`
	Info     string       `json:"info"`
	Income   float64      `json:"income"`
	CaseType string       `json:"case_type"`
}

type QueDetail struct {
	DefaultDuration int                     `json:"default_duration"`
	Limits          []Limit                 `json:"limits"`
	Ques            []SimpleAppointmentInfo `json:"ques"`
	Totals          []TotalLimit            `json:"totals"`
	WorkHours       WorkHour                `json:"work_hours"`
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
