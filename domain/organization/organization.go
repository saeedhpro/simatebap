package organization

import "database/sql"

type OrganizationInfo struct {
	ID           int64                 `json:"id"`
	Name         string                `json:"name"`
	Phone        string                `json:"phone"`
	Phone1       string                `json:"phone1"`
	ProfessionID string                `json:"profession_id"`
	Profession   *SimpleProfessionInfo `json:"profession"`
	KnownAs      string                `json:"known_as"`
	CaseTypes    string                `json:"case_types"`
	Logo         string                `json:"logo"`
	StaffID      int64                 `json:"staff_id"`
	Staff        *SimpleUserInfo       `json:"staff"`
	Info         string                `json:"info"`
	Website      string                `json:"website"`
	Instagram    string                `json:"instagram"`
	SmsCredit    int64                 `json:"sms_credit"`
	SmsPrice     int64                 `json:"sms_price"`
	CreatedAt    *sql.NullTime         `json:"created_at"`
}

type CreateOrganizationRequest struct {
	Name             string                `json:"name" binding:"required"`
	Phone            string                `json:"phone" binding:"required"`
	Phone1           string                `json:"phone1"`
	ProfessionID     int64                 `json:"profession_id" binding:"required"`
	KnownAs          string                `json:"known_as"`
	CaseTypes        string                `json:"case_types"`
	Info             string                `json:"info"`
	Website          string                `json:"website"`
	Instagram        string                `json:"instagram"`
	SmsCredit        int                   `json:"sms_credit"`
	SmsPrice         int                   `json:"sms_price"`
	Logo             string                `json:"logo"`
	RelRadiologies   []RelOrganizationType `json:"rel_radiologies" binding:"required"`
	RelLaboratories  []RelOrganizationType `json:"rel_laboratories" binding:"required"`
	RelDoctorOffices []RelOrganizationType `json:"rel_doctor_offices" binding:"required"`
}

type SetOrganizationSliderRequest struct {
	Images []string `json:"images" binding:"required"`
}

type RelOrganizationType struct {
	ID           int64  `json:"id"`
	ProfessionID string `json:"profession_id"`
}

type UpdateOrganizationRequest struct {
	Name             string                `json:"name"`
	Phone            string                `json:"phone"`
	Phone1           string                `json:"phone1"`
	KnownAs          string                `json:"known_as"`
	CaseTypes        string                `json:"case_types"`
	Info             string                `json:"info"`
	Website          string                `json:"website"`
	Instagram        string                `json:"instagram"`
	SmsCredit        int                   `json:"sms_credit"`
	SmsPrice         int                   `json:"sms_price"`
	Logo             string                `json:"logo"`
	RelRadiologies   []RelOrganizationType `json:"rel_radiologies"`
	RelLaboratories  []RelOrganizationType `json:"rel_laboratories"`
	RelDoctorOffices []RelOrganizationType `json:"rel_doctor_offices"`
}

type UpdateOrganizationWorkHourRequest struct {
	WorkHourStart string `json:"work_hour_start"`
	WorkHourEnd   string `json:"work_hour_end"`
}

type SimpleOrganizationInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SimpleOrganizationVipScheduleInfo struct {
	ID          int64         `json:"id"`
	AppCount    int64         `json:"app_count"`
	SiteCount   int64         `json:"site_count"`
	DoctorCount int64         `json:"doctor_count"`
	StartAt     *sql.NullTime `json:"start_at"`
	EndAt       *sql.NullTime `json:"end_at"`
	CaseType    string        `json:"case_type"`
}

type SimpleVipScheduleCaseInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateOrganizationVipScheduleRequest struct {
	DoctorCount int64  `json:"doctor_count"`
	SiteCount   int64  `json:"site_count"`
	AppCount    int64  `json:"app_count"`
	StartAt     string `json:"start_at" binding:"required"`
	EndAt       string `json:"end_at" binding:"required"`
	CaseType    string `json:"case_type"`
}

type CreateVipScheduleCaseRequest struct {
	Name string `json:"name" binding:"required"`
}

type SimpleProfessionInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProfessionInfo struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	LaboratoryCases  string `json:"laboratory_cases"`
	PhotographyCases string `json:"photography_cases"`
	RadiologyCases   string `json:"radiology_cases"`
}
