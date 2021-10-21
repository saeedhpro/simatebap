package admission

type CreateAdmissionRequest struct {
	AppointmentID string `json:"appointment_id"`
	CaseName      string `json:"case_name"`
	FileName      string `json:"file_name"`
}

type UpdateAdmissionRequest struct {
	Info          string `json:"info"`
	AppointmentID string `json:"appointment_id"`
	CaseName      string `json:"case_name"`
	FileName      string `json:"file_name"`
}

type AdmissionInfo struct {
	ID            int64  `json:"id"`
	AppointmentID string `json:"appointment_id"`
	CaseName      string `json:"case_name"`
	FileName      string `json:"file_name"`
}
