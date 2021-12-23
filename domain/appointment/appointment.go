package appointment

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
	"time"
)

type CreateAppointmentRequest struct {
	UserID   int64   `json:"user_id"`
	Info     string  `json:"info"`
	StartAt  string  `json:"start_at"`
	CaseType string  `json:"case_type"`
	Income   float64 `json:"income"`
	IsVip    int     `json:"is_vip"`
}
type AppointmentSendResultRequest struct {
	Images []Image `json:"images"`
	Type   string  `json:"type"`
}

type Image struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type UpdateAppointmentRequest struct {
	Info               string    `json:"info"`
	StartAt            time.Time `json:"start_at"`
	CaseType           string    `json:"case_type"`
	Income             float64   `json:"income"`
	IsVip              int       `json:"is_vip"`
	FuturePrescription string    `json:"future_prescription"`
	Prescription       string    `json:"prescription"`
	PhotographyCases   []string  `json:"photography_cases"`
	RadiologyCases     []string  `json:"radiology_cases"`
}

type ChangeAppointmentStatusRequest struct {
	Status int `json:"status"`
}

type SimpleAppointmentInfo struct {
	ID        int64                          `json:"id"`
	CaseType  string                         `json:"case_type"`
	Duration  string                         `json:"duration"`
	IsVip     int64                          `json:"is_vip"`
	StartAt   time.Time                      `json:"start_at"`
	Status    int                            `json:"status"`
	UserID    int64                          `json:"user_id"`
	User      *organization.OrganizationUser `json:"user"`
	Info      string                         `json:"info"`
	Income    float64                        `json:"income"`
	Price     float64                        `json:"price"`
	UpdatedAt sql.NullTime                   `json:"updated_at"`
}

type UserAppointmentInfo struct {
	ID                 int64         `json:"id"`
	UserID             int64         `json:"user_id"`
	CreatedAt          *sql.NullTime `json:"created_at"`
	Info               string        `json:"info"`
	StaffID            int64         `json:"staff_id"`
	StartAt            *sql.NullTime `json:"start_at"`
	EndAt              *sql.NullTime `json:"end_at"`
	Status             int           `json:"status"`
	DirectorID         int64         `json:"director_id"`
	UpdatedAt          *sql.NullTime `json:"updated_at"`
	Income             float64       `json:"income"`
	Subject            string        `json:"subject"`
	CaseType           string        `json:"case_type"`
	LaboratoryCases    string        `json:"laboratory_cases"`
	PhotographyStatus   string        `json:"photography_status"`
	PhotographyCases   string        `json:"photography_cases"`
	RadiologyStatus     string        `json:"radiology_status"`
	RadiologyCases     string        `json:"radiology_cases"`
	Prescription       string        `json:"prescription"`
	FuturePrescription string        `json:"future_prescription"`
	LaboratoryMsg      string        `json:"laboratory_msg"`
	PhotographyMsg     string        `json:"photography_msg"`
	RadiologyMsg       string        `json:"radiology_msg"`
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
	LRndImg            string        `json:"l_rnd_img"`
	RRndImg            string        `json:"r_rnd_img"`
	PRndImg            string        `json:"p_rnd_img"`
	LImgs              int           `json:"l_imgs"`
	RImgs              int           `json:"r_imgs"`
	PImgs              int           `json:"p_imgs"`
	Code               string        `json:"code"`
	IsVip              bool          `json:"is_vip"`
	VipIntroducer      int64         `json:"vip_introducer"`
	Absence            int           `json:"absence"`
	FileID             string        `json:"file_id"`
	Price              float64       `json:"price"`
	FName              string        `json:"fname"`
	LName              string        `json:"lname"`
	Tel                string        `json:"tel"`
	Logos              []string      `json:"logos"`
}

type OperationInfo struct {
	ID                int64                          `json:"id"`
	UserID            int64                          `json:"user_id"`
	User              *organization.OrganizationUser `json:"user"`
	StartAt           sql.NullTime                   `json:"start_at"`
	Info              string                         `json:"info"`
	Income            float64                        `json:"income"`
	CaseType          string                         `json:"case_type"`
	Code              string                         `json:"code"`
	Status            int64                          `json:"status"`
	PhotographyStatus int64                          `json:"photography_status"`
	RadiologyStatus   int64                          `json:"radiology_status"`
	FName             string                         `json:"fname"`
	LName             string                         `json:"lname"`
	Tel               string                         `json:"tel"`
	FileID            string                         `json:"file_id"`
	Logo              string                         `json:"logo"`
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

func GetAppointmentById(id string) (*UserAppointmentInfo, error) {
	var appointment UserAppointmentInfo
	query := "SELECT appointment.id id, appointment.user_id user_id, appointment.created_at created_at, ifnull(appointment.info, '') info, appointment.staff_id staff_id, appointment.start_at start_at, appointment.end_at end_at, appointment.status status, ifnull(appointment.director_id, -1) director_id, ifnull(appointment.updated_at, null) updated_at, appointment.income, ifnull(appointment.subject, '') subject, ifnull(appointment.case_type, '') case_type, ifnull(appointment.laboratory_cases, '') laboratory_cases, ifnull(appointment.photography_cases, '') photography_cases, ifnull(appointment.radiology_cases, '') radiology_cases, ifnull(appointment.prescription, '') prescription, ifnull(appointment.future_prescription, '') future_prescription, ifnull(appointment.laboratory_msg, '') laboratory_msg, ifnull(appointment.photography_msg, '') photography_msg, ifnull(appointment.radiology_msg, '') radiology_msg, appointment.organization_id, ifnull(appointment.director_id, -1) laboratory_id, ifnull(appointment.photography_id, -1) photography_id, ifnull(appointment.radiology_id, -1) radiology_id, appointment.l_admission_at, appointment.r_admission_at, appointment.p_admission_at, appointment.l_result_at, appointment.r_result_at, appointment.p_result_at, ifnull(appointment.l_rnd_img, '') l_rnd_img, ifnull(appointment.r_rnd_img, '') r_rnd_img, ifnull(appointment.p_rnd_img, '') p_rnd_img, appointment.l_imgs, appointment.r_imgs, appointment.p_imgs, ifnull(appointment.code, '') code, appointment.is_vip, ifnull(appointment.vip_introducer, 0) vip_introducer, appointment.absence, ifnull(user.file_id, '') file_id, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel FROM appointment LEFT JOIN user on appointment.user_id = user.id WHERE appointment.id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "err")
		return nil, err
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&appointment.ID,
		&appointment.UserID,
		&appointment.CreatedAt,
		&appointment.Info,
		&appointment.StaffID,
		&appointment.StartAt,
		&appointment.EndAt,
		&appointment.Status,
		&appointment.DirectorID,
		&appointment.UpdatedAt,
		&appointment.Income,
		&appointment.Subject,
		&appointment.CaseType,
		&appointment.LaboratoryCases,
		&appointment.PhotographyCases,
		&appointment.RadiologyCases,
		&appointment.Prescription,
		&appointment.FuturePrescription,
		&appointment.LaboratoryMsg,
		&appointment.PhotographyMsg,
		&appointment.RadiologyMsg,
		&appointment.OrganizationID,
		&appointment.LaboratoryID,
		&appointment.PhotographyID,
		&appointment.RadiologyID,
		&appointment.LAdmissionAt,
		&appointment.RAdmissionAt,
		&appointment.PAdmissionAt,
		&appointment.LResultAt,
		&appointment.RResultAt,
		&appointment.PResultAt,
		&appointment.LRndImg,
		&appointment.RRndImg,
		&appointment.PRndImg,
		&appointment.LImgs,
		&appointment.RImgs,
		&appointment.PImgs,
		&appointment.Code,
		&appointment.IsVip,
		&appointment.VipIntroducer,
		&appointment.Absence,
		&appointment.FileID,
		&appointment.FName,
		&appointment.LName,
		&appointment.Tel,
	)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &appointment, nil
}
