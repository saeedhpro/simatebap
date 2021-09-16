package caseType

type CaseTypeInfo struct {
	ID             int64  `json:"id"`
	Name           string `json:"name,omitempty"`
	OrganizationId int64  `json:"organization_id,omitempty"`
	Duration       int    `json:"duration,omitempty"`
	IsLimited      bool   `json:"is_limited,omitempty"`
	Limitation     int    `json:"limitation,omitempty"`
}

type CreateCaseTypeRequest struct {
	Name           string `json:"name,omitempty"`
	OrganizationId int64  `json:"organization_id"`
	Duration       int    `json:"duration"`
	IsLimited      bool   `json:"is_limited"`
	Limitation     int    `json:"limitation"`
}

type UpdateCaseTypeRequest struct {
	Name           string `json:"name,omitempty"`
	OrganizationId int64  `json:"organization_id"`
	Duration       int    `json:"duration"`
	IsLimited      bool   `json:"is_limited"`
	Limitation     int    `json:"limitation"`
}
