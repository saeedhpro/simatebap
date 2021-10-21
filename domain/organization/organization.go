package organization

import "database/sql"

type OrganizationInfo struct {
	ID           int64                `json:"id"`
	Name         string               `json:"name,omitempty"`
	Phone        string               `json:"phone,omitempty"`
	Phone1       string               `json:"phone1,omitempty"`
	ProfessionID string               `json:"profession_id,omitempty"`
	Profession   SimpleProfessionInfo `json:"profession,omitempty"`
	KnownAs      string               `json:"known_as,omitempty"`
	CaseTypes    string               `json:"case_types,omitempty"`
	Logo         string               `json:"logo,omitempty"`
	StaffID      int64                `json:"staff_id"`
	Staff        SimpleUserInfo       `json:"staff"`
	Info         string               `json:"info,omitempty"`
	Website      string               `json:"website,omitempty"`
	Instagram    string               `json:"instagram,omitempty"`
	SmsCredit    int64                `json:"sms_credit"`
	SmsPrice     int64                `json:"sms_price"`
	CreatedAt    sql.NullTime         `json:"created_at"`
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
	ID           int64 `json:"id"`
	ProfessionID int64 `json:"profession_id"`
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

type SimpleOrganizationInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
}

type SimpleProfessionInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
}

type ProfessionInfo struct {
	ID               int64  `json:"id"`
	Name             string `json:"name,omitempty"`
	LaboratoryCases  string `json:"laboratory_cases,omitempty"`
	PhotographyCases string `json:"photography_cases,omitempty"`
	RadiologyCases   string `json:"radiology_cases,omitempty"`
}
