package caseType

type CaseTypeInfo struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	OrganizationId int64  `json:"organization_id"`
	Duration       int    `json:"duration"`
	IsLimited      bool   `json:"is_limited"`
	Limitation     int    `json:"limitation"`
}

type CaseInfo struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ParentID     *int64 `json:"parent_id"`
	IsMain       bool   `json:"is_main"`
	ProfessionID int64  `json:"profession_id"`
}

type CreateCaseTypeRequest struct {
	Name           string `json:"name"`
	OrganizationId int64  `json:"organization_id"`
	Duration       int    `json:"duration"`
	IsLimited      bool   `json:"is_limited"`
	Limitation     int    `json:"limitation"`
}

type UpdateCaseTypeRequest struct {
	Name           string `json:"name"`
	OrganizationId int64  `json:"organization_id"`
	Duration       int    `json:"duration"`
	IsLimited      bool   `json:"is_limited"`
	Limitation     int    `json:"limitation"`
}
