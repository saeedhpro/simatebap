package organizationController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	appointment2 "gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	wallet2 "gitlab.com/simateb-project/simateb-backend/domain/wallet"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/repository/vip"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"gitlab.com/simateb-project/simateb-backend/utils/pagination"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type OrganizationControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	GetList(c *gin.Context)
	GetListAll(c *gin.Context)
	GetOrganizationAppointments(c *gin.Context)
	UploadOrganizationImage(c *gin.Context)
	UpdateOrganizationAbout(c *gin.Context)
	GetOrganizationImages(c *gin.Context)
	GetOrganizationAbout(c *gin.Context)
	GetOrganizationWorkTime(c *gin.Context)
	UpdateOrganizationWorkTime(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	GetUsers(c *gin.Context)
	GetEmployees(c *gin.Context)
	GetOrganizationWallet(c *gin.Context)
	IncreaseOrganizationWallet(c *gin.Context)
	DecreaseOrganizationWallet(c *gin.Context)
	SetOrganizationWallet(c *gin.Context)
	GetOrganizationRelList(c *gin.Context)
	GetOrganizationRelOfficesList(c *gin.Context)
	SetOrganizationSlider(c *gin.Context)
	GetOrganizationScheduleList(c *gin.Context)
	GetOrganizationScheduleCasesList(c *gin.Context)
	GetOrganizationSchedule(c *gin.Context)
	CreateOrganizationSchedule(c *gin.Context)
	CreateOrganizationScheduleCase(c *gin.Context)
	GetVipScheduleCase(c *gin.Context)
}

type OrganizationControllerStruct struct {
}

func NewOrganizationController() OrganizationControllerInterface {
	x := &OrganizationControllerStruct{
	}
	return x
}

func (oc *OrganizationControllerStruct) Create(c *gin.Context) {
	var createOrganizationRequest organization.CreateOrganizationRequest
	if errors := c.ShouldBindJSON(&createOrganizationRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateOrganizationQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	staffID := auth.GetStaffUser(c).UserID
	result, err := stmt.Exec(
		createOrganizationRequest.Name,
		createOrganizationRequest.KnownAs,
		createOrganizationRequest.ProfessionID,
		createOrganizationRequest.Logo,
		createOrganizationRequest.Phone,
		createOrganizationRequest.Phone1,
		staffID,
		createOrganizationRequest.Info,
		createOrganizationRequest.CaseTypes,
		createOrganizationRequest.SmsPrice,
		createOrganizationRequest.SmsCredit,
		createOrganizationRequest.Website,
		createOrganizationRequest.Instagram,
	)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	setOrganizationRelations(id, createOrganizationRequest.RelRadiologies, createOrganizationRequest.RelLaboratories, createOrganizationRequest.RelDoctorOffices)
	c.JSON(http.StatusOK, true)
}

func setOrganizationRelations(id int64, radiologies []organization.RelOrganizationType, laboratories []organization.RelOrganizationType, offices []organization.RelOrganizationType) {
	var ids []int64
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetOrganizationRelations)
	if err != nil {
		return
	}
	rows, error := stmt.Query(id)
	if error != nil {
		return
	}
	for rows.Next() {
		var rel_organization_id int64
		rows.Scan(&rel_organization_id)
		ids = append(ids, rel_organization_id)
	}
	insertQuery := "INSERT INTO `rel_organization`(`organization_id`, `profession_id`, `rel_organization_id`) VALUES "
	var radValues []organization.RelOrganizationType
	var labValues []organization.RelOrganizationType
	var offValues []organization.RelOrganizationType
	var queryStr []string
	var allValues []interface{}
	for _, n := range radiologies {
		if exists := helper.ItemExists(ids, n); !exists {
			radValues = append(radValues, n)
		}
	}
	for _, n := range laboratories {
		if !helper.ItemExists(ids, n) {
			labValues = append(labValues, n)
		}
	}
	for _, n := range offices {
		if !helper.ItemExists(ids, n) {
			offValues = append(offValues, n)
		}
	}
	for _, i := range radValues {
		queryStr = append(queryStr, "(?,?,?)")
		allValues = append(allValues, id, i.ProfessionID, i.ID)
	}
	for _, i := range labValues {
		queryStr = append(queryStr, "(?,?,?)")
		allValues = append(allValues, id, i.ProfessionID, i.ID)
	}
	for _, i := range offValues {
		queryStr = append(queryStr, "(?,?,?)")
		allValues = append(allValues, id, i.ProfessionID, i.ID)
	}
	insertQuery = fmt.Sprintf("%s%s", insertQuery, strings.Join(queryStr, ","))
	stmt, error = repository.DBS.MysqlDb.Prepare(insertQuery)
	if error != nil {
		log.Println(error.Error())
		return
	}
	_, error = stmt.Exec(allValues...)
	if error != nil {
		log.Println(error.Error())
	}
}

func (oc *OrganizationControllerStruct) Get(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		return
	}
	GetOrganizationQuery := "SELECT id, ifnull(name, ''), ifnull(phone, ''), ifnull(phone1, ''), ifnull(profession_id, ''), ifnull(known_as, ''), ifnull(case_types, ''), ifnull(staff_id, ''), ifnull(info, ''), ifnull(website, ''), ifnull(instagram, ''), sms_price, sms_credit, created_at, ifnull(logo, '') logo FROM organization WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(GetOrganizationQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var organizationInfo organization.OrganizationInfo
	result := stmt.QueryRow(organizationID)
	err = result.Scan(
		&organizationInfo.ID,
		&organizationInfo.Name,
		&organizationInfo.Phone,
		&organizationInfo.Phone1,
		&organizationInfo.ProfessionID,
		&organizationInfo.KnownAs,
		&organizationInfo.CaseTypes,
		&organizationInfo.StaffID,
		&organizationInfo.Info,
		&organizationInfo.Website,
		&organizationInfo.Instagram,
		&organizationInfo.SmsPrice,
		&organizationInfo.SmsCredit,
		&organizationInfo.CreatedAt,
		&organizationInfo.Logo,
	)
	organizationInfo.Profession = GetProfession(organizationInfo.ProfessionID)
	organizationInfo.Staff = getStaff(organizationInfo.StaffID)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "یافت نشد")
			return
		}
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, organizationInfo)
}

func (oc *OrganizationControllerStruct) GetList(c *gin.Context) {
	var query = "SELECT id, ifnull(name, ''), ifnull(phone, ''), ifnull(phone1, ''), ifnull(profession_id, '')," +
		" ifnull(known_as, ''), ifnull(case_types, ''), ifnull(staff_id, ''), ifnull(info, ''), ifnull(website, '')," +
		" ifnull(instagram, ''), sms_price, sms_credit, created_at, ifnull(logo, '') logo FROM organization "
	var values []interface{}
	userGroupID := c.Query("group")
	q := c.Query("q")
	var query2 = ""
	if q != "" && q != "null" && q != "undefined" {
		q = "'%" + q + "%'"
		query2 += fmt.Sprintf(" WHERE name LIKE %s ", q)
	}
	var err error
	if userGroupID != "" {
		if userGroupID == "2" || userGroupID == "3" {
			if query2 != "" {
				query += query2
				query += " AND profession_id = ? "
			} else {
				query += " WHERE profession_id = ? "
			}
			values = append(values, userGroupID)
		} else {
			if query2 != "" {
				query += query2
				query += " AND profession_id != 2 AND profession_id != 3 "
			} else {
				query += " WHERE profession_id != 2 AND profession_id != 3 "
			}
		}
	} else {
		query += query2
	}
	query += " ORDER BY id DESC "

	page := c.Query("page")
	if page != "" && page != "null" && page != "undefined" {
		offset, _ := strconv.Atoi(page)
		offset = (offset - 1) * 10
		values = append(values, offset)
		query += " LIMIT 10 OFFSET ?"
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "log")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	organizations := []organization.OrganizationInfo{}
	var organizationInfo organization.OrganizationInfo
	rows, error := stmt.Query(values...)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&organizationInfo.ID,
			&organizationInfo.Name,
			&organizationInfo.Phone,
			&organizationInfo.Phone1,
			&organizationInfo.ProfessionID,
			&organizationInfo.KnownAs,
			&organizationInfo.CaseTypes,
			&organizationInfo.StaffID,
			&organizationInfo.Info,
			&organizationInfo.Website,
			&organizationInfo.Instagram,
			&organizationInfo.SmsPrice,
			&organizationInfo.SmsCredit,
			&organizationInfo.CreatedAt,
			&organizationInfo.Logo,
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
		profession := GetProfession(organizationInfo.ProfessionID)
		if profession != nil {
			organizationInfo.Profession = profession
		}
		staff := getStaff(organizationInfo.StaffID)
		if staff != nil {
			organizationInfo.Staff = staff
		}
		organizations = append(organizations, organizationInfo)
	}
	p, err := strconv.Atoi(page)
	count := 0
	count, _ = getOrganizationCountAdmin(q, userGroupID)
	paginated := pagination.OrganizationPaginationInfo{
		Data:       organizations,
		Page:       p,
		PagesCount: count,
	}
	if p > 1 {
		paginated.PrevPage = p - 1
	} else {
		paginated.PrevPage = p
	}
	if p < count/10 {
		paginated.NextPage = p
	} else {
		paginated.NextPage = p + 1
	}
	paginated.HasNextPage = (bool)(count > 10 && count > (p*10))
	c.JSON(http.StatusOK, paginated)
}

func getOrganizationCountAdmin(q string, userGroupID string) (int, error) {
	query := "SELECT COUNT(*) FROM organization "
	var values []interface{}
	count := 0
	var query2 = ""
	if q != "" && q != "null" && q != "undefined" {
		query2 += fmt.Sprintf(" WHERE name LIKE %s ", q)
	}
	if userGroupID != "" {
		if userGroupID == "2" || userGroupID == "3" {
			if query2 != "" {
				query += query2
				query += fmt.Sprintf(" AND profession_id = ? ")
			} else {
				query += fmt.Sprintf(" WHERE profession_id = ? ")
			}
			values = append(values, userGroupID)
		} else {
			if query2 != "" {
				query += query2
				query += fmt.Sprintf(" AND profession_id != 2 AND profession_id != 3 ")
			} else {
				query += fmt.Sprintf(" WHERE profession_id != 2 AND profession_id != 3 ")
			}
		}
	} else {
		query += query2
	}
	query += " ORDER BY id DESC "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func (oc *OrganizationControllerStruct) GetListAll(c *gin.Context) {
	var query = "SELECT id, ifnull(name, ''), ifnull(phone, ''), ifnull(phone1, ''), ifnull(profession_id, '')," +
		" ifnull(known_as, ''), ifnull(case_types, ''), ifnull(staff_id, ''), ifnull(info, ''), ifnull(website, '')," +
		" ifnull(instagram, ''), sms_price, sms_credit, created_at FROM organization "
	var values []interface{}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "log")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	organizations := []organization.OrganizationInfo{}
	var organizationInfo organization.OrganizationInfo
	rows, error := stmt.Query(values...)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&organizationInfo.ID,
			&organizationInfo.Name,
			&organizationInfo.Phone,
			&organizationInfo.Phone1,
			&organizationInfo.ProfessionID,
			&organizationInfo.KnownAs,
			&organizationInfo.CaseTypes,
			&organizationInfo.StaffID,
			&organizationInfo.Info,
			&organizationInfo.Website,
			&organizationInfo.Instagram,
			&organizationInfo.SmsPrice,
			&organizationInfo.SmsCredit,
			&organizationInfo.CreatedAt,
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
		profession := GetProfession(organizationInfo.ProfessionID)
		if profession != nil {
			organizationInfo.Profession = profession
		}
		staff := getStaff(organizationInfo.StaffID)
		if staff != nil {
			organizationInfo.Staff = staff
		}
		organizations = append(organizations, organizationInfo)
	}
	c.JSON(http.StatusOK, organizations)
}

func (oc *OrganizationControllerStruct) UploadOrganizationImage(c *gin.Context) {
	orgID := c.Param("id")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err.Error(), "read file")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	t := time.Now().UnixNano()
	fileName := fmt.Sprintf("%d%s", t, filepath.Ext(header.Filename))
	path := fmt.Sprintf("./images/organizations/%s", orgID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Println(err.Error())
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
	}
	err = c.SaveUploadedFile(header, path+"/"+fileName)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"name": fileName,
		"id":   orgID,
		"path": fmt.Sprintf("/images/organizations/%s/%s", orgID, fileName),
	})
}

func (oc *OrganizationControllerStruct) UpdateOrganizationAbout(c *gin.Context) {
	orgID := c.Param("id")
	var values []interface{}
	var request pagination.OrganizationAboutRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	var query = "UPDATE `organization` SET `text1`=?,`image1`=?,`text2`=?,`image2`=?,`text3`=?,`image3`=?,`text4`=?,`image4`=? WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(
		values,
		request.Text1,
		request.Image1,
		request.Text2,
		request.Image2,
		request.Text3,
		request.Image3,
		request.Text4,
		request.Image4,
		orgID,
	)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func (oc *OrganizationControllerStruct) GetOrganizationImages(c *gin.Context) {
	id := c.Param("id")
	logos := []string{}
	files, err := ioutil.ReadDir(fmt.Sprintf("./images/results/%s", id))
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() != "." || f.Name() != ".." {
			logos = append(logos, fmt.Sprintf("http://%s/images/organizations/%s/%s", c.Request.Host, id, f.Name()))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"logos": logos,
	})
}

func (oc *OrganizationControllerStruct) GetOrganizationAbout(c *gin.Context) {
	query := "SELECT ifnull(text1, '') text1, ifnull(text2, '') text2, ifnull(text3, '') text3, ifnull(text4, '') text4, ifnull(image1, '') image1, ifnull(image2, '') image2, ifnull(image3, '') image3, ifnull(image4, '') image4 FROM organization WHERE id = ? "
	orgID := c.Param("id")
	var values []interface{}
	values = append(values, orgID)
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	image := pagination.OrganizationAbout{}
	if err != nil {
		log.Println(err.Error(), "prepare")
		c.JSON(200, image)
		return
	}
	row := stmt.QueryRow(values...)
	err = row.Scan(
		&image.Text1,
		&image.Text2,
		&image.Text3,
		&image.Text4,
		&image.Image1,
		&image.Image2,
		&image.Image3,
		&image.Image4,
	)
	c.JSON(http.StatusOK, image)
}

func (oc *OrganizationControllerStruct) GetOrganizationAppointments(c *gin.Context) {
	query := "SELECT appointment.id id, appointment.user_id user_id, appointment.created_at created_at, ifnull(appointment.info, '') info, appointment.staff_id staff_id, appointment.start_at start_at, appointment.end_at end_at, appointment.status status, ifnull(appointment.director_id, -1) director_id, ifnull(appointment.updated_at, null) updated_at, appointment.income, ifnull(appointment.subject, '') subject, ifnull(appointment.case_type, '') case_type, ifnull(appointment.laboratory_cases, '') laboratory_cases, ifnull(appointment.photography_cases, '') photography_cases, ifnull(appointment.radiology_cases, '') radiology_cases, ifnull(appointment.prescription, '') prescription, ifnull(appointment.future_prescription, '') future_prescription, ifnull(appointment.laboratory_msg, '') laboratory_msg, ifnull(appointment.photography_msg, '') photography_msg, ifnull(appointment.radiology_msg, '') radiology_msg, appointment.organization_id, ifnull(appointment.director_id, -1) laboratory_id, ifnull(appointment.photography_id, -1) photography_id, ifnull(appointment.radiology_id, -1) radiology_id, appointment.l_admission_at, appointment.r_admission_at, appointment.p_admission_at, appointment.l_result_at, appointment.r_result_at, appointment.p_result_at, ifnull(appointment.l_rnd_img, '') l_rnd_img, ifnull(appointment.r_rnd_img, '') r_rnd_img, ifnull(appointment.p_rnd_img, '') p_rnd_img, appointment.l_imgs, appointment.r_imgs, appointment.p_imgs, ifnull(appointment.code, '') code, appointment.is_vip, ifnull(appointment.vip_introducer, 0) vip_introducer, appointment.absence, ifnull(user.file_id, '') file_id, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel FROM appointment LEFT JOIN user on appointment.user_id = user.id WHERE appointment.organization_id = ? "
	orgID := c.Param("id")
	page := c.Query("page")
	q := c.Query("q")
	startAt := c.Query("start_at")
	endAt := c.Query("end_at")
	var values []interface{}
	values = append(values, orgID)
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%' ) "
	}
	if startAt != "" && startAt != "null" && startAt != "undefined" {
		query += " AND start_at >= ?"
		values = append(values, startAt)
	}
	if endAt != "" && endAt != "null" && endAt != "undefined" {
		query += " AND start_at <= ?"
		values = append(values, endAt)
	}
	if page != "" && page != "null" && page != "undefined" {
		query += " LIMIT 10 OFFSET ?"
		values = append(values, orgID)
	}
	log.Println()
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	paginated := pagination.OrganizationAppointmentPaginationInfo{}
	if err != nil {
		log.Println(err.Error(), "prepare")
		c.JSON(200, paginated)
		return
	}
	appointments := []appointment2.UserAppointmentInfo{}
	var appointment appointment2.UserAppointmentInfo
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, paginated)
		return
	}
	for rows.Next() {
		err = rows.Scan(
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
			log.Println(err.Error(), "err")
			c.JSON(http.StatusOK, paginated)
			return
		}
		appointments = append(appointments, appointment)
	}
	p, err := strconv.Atoi(page)
	count := 0
	count, _ = GetOrgApponints(orgID, q)
	paginated = pagination.OrganizationAppointmentPaginationInfo{
		Data:       appointments,
		Page:       p,
		PagesCount: count,
	}
	c.JSON(http.StatusOK, paginated)
}

func (oc *OrganizationControllerStruct) GetOrganizationWorkTime(c *gin.Context) {
	query := "SELECT ifnull(work_hour_start, '15:00:00') work_hour_start, ifnull(work_hour_end, '21:00:00') work_hour_end FROM Organization WHERE id = ? "
	orgID := c.Param("id")
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	wh := pagination.OrganizationWorkTimeStruct{}
	if err != nil {
		log.Println(err.Error(), "prepare")
		c.JSON(200, wh)
		return
	}
	row := stmt.QueryRow(orgID)
	err = row.Scan(
		&wh.WorkHourStart,
		&wh.WorkHourEnd,
	)
	if err != nil {
		log.Println(err.Error(), "err")
		c.JSON(http.StatusOK, wh)
		return
	}
	c.JSON(http.StatusOK, wh)
}

func (oc *OrganizationControllerStruct) UpdateOrganizationWorkTime(c *gin.Context) {
	query := "UPDATE `organization` SET `work_hour_start`= ? ,`work_hour_end`= ? WHERE `id` = ?"
	orgID := c.Param("id")
	var request organization.UpdateOrganizationWorkHourRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	_, err = stmt.Exec(request.WorkHourStart, request.WorkHourEnd, orgID)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(200, true)
}

func GetOrgApponints(orgID string, q string) (int, error) {
	count := 0
	query := "SELECT COUNT(*) FROM appointment WHERE organization_id = ? "
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%' ) "
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "count error")
		return 0, err
	}
	result := stmt.QueryRow(orgID)
	if err = result.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (oc *OrganizationControllerStruct) GetListForAdmin(c *gin.Context) {
	var query = "SELECT id, ifnull(name, ''), ifnull(phone, ''), ifnull(phone1, ''), ifnull(profession_id, '')," +
		" ifnull(known_as, ''), ifnull(case_types, ''), ifnull(staff_id, ''), ifnull(info, ''), ifnull(website, '')," +
		" ifnull(instagram, ''), sms_price, sms_credit, created_at FROM organization "
	var values []interface{}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		q = "'%" + q + "%'"
		query += fmt.Sprintf(" WHERE name LIKE %s ", q)
	}
	var err error
	query += " ORDER BY id DESC "
	page := c.Query("page")
	if page != "" && page != "null" && page != "undefined" {
		offset, _ := strconv.Atoi(page)
		offset = (offset - 1) * 10
		values = append(values, offset)
		query += " LIMIT 10 OFFSET ?"
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "log")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	organizations := []organization.OrganizationInfo{}
	var organizationInfo organization.OrganizationInfo
	rows, error := stmt.Query(values...)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&organizationInfo.ID,
			&organizationInfo.Name,
			&organizationInfo.Phone,
			&organizationInfo.Phone1,
			&organizationInfo.ProfessionID,
			&organizationInfo.KnownAs,
			&organizationInfo.CaseTypes,
			&organizationInfo.StaffID,
			&organizationInfo.Info,
			&organizationInfo.Website,
			&organizationInfo.Instagram,
			&organizationInfo.SmsPrice,
			&organizationInfo.SmsCredit,
			&organizationInfo.CreatedAt,
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
		profession := GetProfession(organizationInfo.ProfessionID)
		if profession != nil {
			organizationInfo.Profession = profession
		}
		staff := getStaff(organizationInfo.StaffID)
		if staff != nil {
			organizationInfo.Staff = staff
		}
		organizations = append(organizations, organizationInfo)
	}
	i, err := strconv.Atoi(page)
	paginated := pagination.OrganizationPaginationInfo{
		Data:        organizations,
		HasNextPage: true,
		PrevPage:    -1,
		NextPage:    2,
		Page:        i,
	}
	c.JSON(http.StatusOK, paginated)
}

func GetProfession(id string) *organization.SimpleProfessionInfo {
	query := "SELECT id, ifnull(name, '') name FROM profession WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	var professionInfo organization.SimpleProfessionInfo
	if err != nil {
		log.Println(err.Error(), "err get prof")
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&professionInfo.ID,
		&professionInfo.Name,
	)
	if err != nil {
		log.Println(err.Error(), "err get prof")
		return nil
	}
	return &professionInfo
}

func getStaff(id int64) *organization.SimpleUserInfo {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetSimpleStaffQuery)
	var userInfo organization.SimpleUserInfo
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&userInfo.ID,
		&userInfo.FirstName,
		&userInfo.LastName,
		&userInfo.Organization,
	)
	if err != nil {
		return nil
	}
	return &userInfo
}

func (oc *OrganizationControllerStruct) GetUsers(c *gin.Context) {
	var values []interface{}
	organizationID := c.Param("id")
	values = append(values, organizationID)
	userGroupID := c.Query("group")
	if userGroupID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد group الزامی است",
			Code:    422,
		})
		return
	}
	values = append(values, userGroupID)
	query := "SELECT user.id id, ifnull(user.nid, '') nid, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel, ifnull(user.user_group_id, '') user_group_id, user.created created, user.last_login last_login, user.birth_date birth_date, user_group.name user_group_name, ifnull(organization.name, '') organization_name, organization.id organization_id, ifnull(user.file_id, '') file_id FROM (user LEFT JOIN user_group on user.user_group_id = user_group.id) LEFT JOIN organization on organization.id = user.organization_id WHERE user.organization_id = ? and user.user_group_id = ? "
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" {
		query += " AND user.created >= ? "
		values = append(values, startDate)
	}
	if endDate != "" {
		query += " AND user.created <= ? "
		values = append(values, endDate)
	}
	query += " ORDER BY id DESC"
	page := c.Query("page")
	if page != "" {
		offset, _ := strconv.Atoi(page)
		offset = (offset - 1) * 10
		values = append(values, offset)
		query += " LIMIT 10 OFFSET ?"
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	paginationInfo := pagination.OrganizationUserPaginationInfo{}
	users := []organization.OrganizationUser{}
	var user organization.OrganizationUser
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Nid,
			&user.FirstName,
			&user.LastName,
			&user.Tel,
			&user.UserGroupID,
			&user.Created,
			&user.LastLogin,
			&user.BirthDate,
			&user.UserGroupName,
			&user.OrganizationName,
			&user.OrganizationID,
			&user.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "user log")
			return
		}
		users = append(users, user)
	}
	paginationInfo.Data = users
	count := 0
	p, _ := strconv.Atoi(page)
	count, err = getOrganizationUsersCount(organizationID, userGroupID, q, startDate, endDate)
	paginationInfo.PagesCount = count
	paginationInfo.Page = p
	if p > 1 {
		paginationInfo.PrevPage = p - 1
	} else {
		paginationInfo.PrevPage = p
	}
	if p < count/10 {
		paginationInfo.NextPage = p
	} else {
		paginationInfo.NextPage = p + 1
	}
	paginationInfo.HasNextPage = (bool)(count > 10 && count > (p*10))
	c.JSON(http.StatusOK, paginationInfo)
}

func (oc *OrganizationControllerStruct) GetEmployees(c *gin.Context) {
	var values []interface{}
	organizationID := c.Param("id")
	values = append(values, organizationID)
	query := "SELECT user.id id, ifnull(user.nid, '') nid, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel, ifnull(user.user_group_id, '') user_group_id, user.created created, user.last_login last_login, user.birth_date birth_date, user_group.name user_group_name, ifnull(organization.name, '') organization_name, organization.id organization_id, ifnull(user.file_id, '') file_id FROM (user LEFT JOIN user_group on user.user_group_id = user_group.id) LEFT JOIN organization on organization.id = user.organization_id WHERE user.organization_id = ? and user.user_group_id NOT IN (1) "
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" {
		query += " AND user.created >= ? "
		values = append(values, startDate)
	}
	if endDate != "" {
		query += " AND user.created <= ? "
		values = append(values, endDate)
	}
	page := c.Query("page")
	if page != "" {
		offset, _ := strconv.Atoi(page)
		offset = (offset - 1) * 10
		values = append(values, offset)
		query += " LIMIT 10 OFFSET ?"
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	paginationInfo := pagination.OrganizationUserPaginationInfo{}
	users := []organization.OrganizationUser{}
	var user organization.OrganizationUser
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Nid,
			&user.FirstName,
			&user.LastName,
			&user.Tel,
			&user.UserGroupID,
			&user.Created,
			&user.LastLogin,
			&user.BirthDate,
			&user.UserGroupName,
			&user.OrganizationName,
			&user.OrganizationID,
			&user.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "user log")
			return
		}
		users = append(users, user)
	}
	paginationInfo.Data = users
	count := 0
	p, _ := strconv.Atoi(page)
	count, err = getOrganizationEmployeeCount(organizationID, q, startDate, endDate)
	paginationInfo.PagesCount = count
	paginationInfo.Page = p
	if p > 1 {
		paginationInfo.PrevPage = p - 1
	} else {
		paginationInfo.PrevPage = p
	}
	if p < count/10 {
		paginationInfo.NextPage = p
	} else {
		paginationInfo.NextPage = p + 1
	}
	paginationInfo.HasNextPage = (bool)(count > 10 && count > (p*10))
	c.JSON(http.StatusOK, paginationInfo)
}

func getOrganizationUsersCount(organizationID string, userGroupID string, q string, startDate string, endDate string) (int, error) {
	query := "SELECT COUNT(*) FROM (user LEFT JOIN user_group on user.user_group_id = user_group.id) LEFT JOIN organization on organization.id = user.organization_id WHERE user.organization_id = ? and user.user_group_id = ? "
	var values []interface{}
	values = append(values, organizationID)
	values = append(values, userGroupID)
	count := 0
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	if startDate != "" {
		query += " AND user.created >= ? "
		values = append(values, startDate)
	}
	if endDate != "" {
		query += " AND user.created <= ? "
		values = append(values, endDate)
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func getOrganizationEmployeeCount(organizationID string, q string, startDate string, endDate string) (int, error) {
	query := "SELECT COUNT(*) FROM (user LEFT JOIN user_group on user.user_group_id = user_group.id) LEFT JOIN organization on organization.id = user.organization_id WHERE user.organization_id = ? and user.user_group_id NOT IN (1) "
	var values []interface{}
	values = append(values, organizationID)
	count := 0
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	if startDate != "" {
		query += " AND user.created >= ? "
		values = append(values, startDate)
	}
	if endDate != "" {
		query += " AND user.created <= ? "
		values = append(values, endDate)
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func (oc *OrganizationControllerStruct) GetOrganizationRelList(c *gin.Context) {
	organizationID := c.Param("id")
	professionID := c.Query("prof")
	if organizationID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد id الزامی است",
			Code:    422,
		})
		return
	}
	if professionID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد prof الزامی است",
			Code:    422,
		})
		return
	}
	getOrganizationRelQuery := "SELECT `organization`.`id` id, ifnull(`organization`.`name`, '') organization_name FROM `organization` LEFT JOIN `rel_organization` ON `organization`.`id` = `rel_organization`.`rel_organization_id` WHERE `rel_organization`.`organization_id` = ? AND `organization`.`profession_id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(getOrganizationRelQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	orgs := []organization.SimpleOrganizationInfo{}
	var org organization.SimpleOrganizationInfo
	rows, err := stmt.Query(organizationID, professionID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&org.ID,
			&org.Name,
		)
		if err != nil {
			log.Println(err.Error(), "Organization")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		orgs = append(orgs, org)
	}
	c.JSON(http.StatusOK, orgs)
}

func (oc *OrganizationControllerStruct) GetOrganizationRelOfficesList(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد id الزامی است",
			Code:    422,
		})
		return
	}
	getOrganizationRelQuery := "SELECT `organization`.`id` id, `organization`.`name` organization_name FROM `organization` LEFT JOIN `rel_organization` ON `organization`.`id` = `rel_organization`.`rel_organization_id` WHERE `rel_organization`.`organization_id` = ? AND `organization`.`profession_id` NOT IN (1,2,3)"
	stmt, err := repository.DBS.MysqlDb.Prepare(getOrganizationRelQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	orgs := []organization.SimpleOrganizationInfo{}
	var org organization.SimpleOrganizationInfo
	rows, err := stmt.Query(organizationID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&org.ID,
			&org.Name,
		)
		if err != nil {
			log.Println(err.Error(), "Organization")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		orgs = append(orgs, org)
	}
	c.JSON(http.StatusOK, orgs)
}

func (oc *OrganizationControllerStruct) GetOrganizationScheduleList(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد id الزامی است",
			Code:    422,
		})
		return
	}
	getOrganizationScheduleQuery := "SELECT id, doctor_count, start_at, end_at FROM `vip_schedule` WHERE organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(getOrganizationScheduleQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	vips := []organization.SimpleOrganizationVipScheduleInfo{}
	var vip organization.SimpleOrganizationVipScheduleInfo
	rows, err := stmt.Query(organizationID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&vip.ID,
			&vip.DoctorCount,
			&vip.StartAt,
			&vip.EndAt,
		)
		if err != nil {
			log.Println(err.Error(), "Organization Vips")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		vips = append(vips, vip)
	}
	c.JSON(http.StatusOK, vips)
}

func (oc *OrganizationControllerStruct) GetOrganizationScheduleCasesList(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد id الزامی است",
			Code:    422,
		})
		return
	}
	getOrganizationScheduleQuery := "SELECT id, name FROM `vip_cases` WHERE organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(getOrganizationScheduleQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	vips := []organization.SimpleVipScheduleCaseInfo{}
	var vip organization.SimpleVipScheduleCaseInfo
	rows, err := stmt.Query(organizationID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&vip.ID,
			&vip.Name,
		)
		if err != nil {
			log.Println(err.Error(), "Organization Vips")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		vips = append(vips, vip)
	}
	c.JSON(http.StatusOK, vips)
}

func (oc *OrganizationControllerStruct) GetOrganizationSchedule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد vip الزامی است",
			Code:    422,
		})
		return
	}
	vid, _ := strconv.ParseInt(id, 10, 64)
	v, _ := vip.GetVipScheduleByID(vid)
	c.JSON(http.StatusOK, v)
}

func (oc *OrganizationControllerStruct) GetVipScheduleCase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد vip الزامی است",
			Code:    422,
		})
		return
	}
	vid, _ := strconv.ParseInt(id, 10, 64)
	v, _ := vip.GetVipScheduleCaseByID(vid)
	c.JSON(http.StatusOK, v)
}

func (oc *OrganizationControllerStruct) Update(c *gin.Context) {
	var updateOrganizationQuery = "UPDATE `organization` SET"
	var values []interface{}
	var columns []string
	organizationId := c.Param("id")
	if organizationId == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var updateOrganizationRequest organization.UpdateOrganizationRequest
	if errors := c.ShouldBindJSON(&updateOrganizationRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getOrganizationUpdateColumns(&updateOrganizationRequest, &columns, &values)
	columnsString := strings.Join(columns, ",")
	updateOrganizationQuery += columnsString
	updateOrganizationQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateOrganizationQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, organizationId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	id, err := strconv.ParseInt(organizationId, 10, 64)
	if err == nil {
		setOrganizationRelations(id, updateOrganizationRequest.RelRadiologies, updateOrganizationRequest.RelLaboratories, updateOrganizationRequest.RelDoctorOffices)
	}
	c.JSON(200, true)
}

func getOrganizationUpdateColumns(o *organization.UpdateOrganizationRequest, columns *[]string, values *[]interface{}) {
	if o.Name != "" {
		*columns = append(*columns, " `name` = ? ")
		*values = append(*values, o.Name)
	}
	if o.Phone != "" {
		*columns = append(*columns, " `phone` = ? ")
		*values = append(*values, o.Phone)
	}
	if o.Phone1 != "" {
		*columns = append(*columns, " `phone1` = ? ")
		*values = append(*values, o.Phone1)
	}
	if o.KnownAs != "" {
		*columns = append(*columns, " `known_as` = ? ")
		*values = append(*values, o.KnownAs)
	}
	if o.CaseTypes != "" {
		*columns = append(*columns, " `case_types` = ? ")
		*values = append(*values, o.CaseTypes)
	}
	if o.Info != "" {
		*columns = append(*columns, " `info` = ? ")
		*values = append(*values, o.Info)
	}
	if o.Website != "" {
		*columns = append(*columns, " `website` = ? ")
		*values = append(*values, o.Website)
	}
	if o.Instagram != "" {
		*columns = append(*columns, " `instagram` = ? ")
		*values = append(*values, o.Instagram)
	}
	*columns = append(*columns, " `sms_credit` = ? ")
	*values = append(*values, o.SmsCredit)
	*columns = append(*columns, " `sms_price` = ? ")
	*values = append(*values, o.SmsPrice)
	if o.Logo != "" {
		*columns = append(*columns, " `logo` = ? ")
		*values = append(*values, o.Logo)
	}
}

func GetOrganization(id string) *organization.OrganizationInfo {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetSimpleOrganizationQuery)
	var organizationInfo organization.OrganizationInfo
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&organizationInfo.ID,
		&organizationInfo.Name,
		&organizationInfo.ProfessionID,
	)
	if err != nil {
		return nil
	}
	return &organizationInfo
}

func GetGroup(id int64) *organization.UserGroup {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetUserGroupQuery)
	var userGroup organization.UserGroup
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(&userGroup.ID, &userGroup.Name)
	if err != nil {
		return nil
	}
	return &userGroup
}

func (oc *OrganizationControllerStruct) GetOrganizationWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(500, gin.H{
			"message": "آی دی صحیح نیست",
		})
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	wallet := wallet2.GetWallet(uID, "organization")
	c.JSON(200, wallet)
}

func (oc *OrganizationControllerStruct) IncreaseOrganizationWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request wallet2.ChangeUserWalletBalance
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, nil)
		return
	}
	wallet := wallet2.GetWallet(uID, "organization")
	result, balance := wallet.Increase(request.Amount)
	if result {
		c.JSON(200, balance)
		return
	}
	c.JSON(500, nil)
}

func (oc *OrganizationControllerStruct) DecreaseOrganizationWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request wallet2.ChangeUserWalletBalance
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, nil)
		return
	}
	wallet := wallet2.GetWallet(uID, "organization")
	result, balance := wallet.Decrease(request.Amount, false)
	if result {
		c.JSON(200, balance)
		return
	}
	c.JSON(500, nil)
}

func (oc *OrganizationControllerStruct) SetOrganizationWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request wallet2.ChangeUserWalletBalance
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, nil)
		return
	}
	wallet := wallet2.GetWallet(uID, "organization")
	result := wallet.SetBalance(request.Amount)
	if result {
		c.JSON(200, nil)
		return
	}
	c.JSON(500, nil)
}

func (oc *OrganizationControllerStruct) SetOrganizationSlider(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		return
	}
	var request organization.SetOrganizationSliderRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	var values []interface{}
	columns := []string{}
	values = append(values, strings.Join(request.Images, ","))
	columns = append(columns, "sliders")
	columnsString := strings.Join(columns, ",")
	var updateOrganizationQuery = "UPDATE `organization` SET"
	updateOrganizationQuery += columnsString
	updateOrganizationQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateOrganizationQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, organizationID)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func (oc *OrganizationControllerStruct) CreateOrganizationSchedule(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		return
	}
	var request organization.CreateOrganizationVipScheduleRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	var values []interface{}
	values = append(values, request.DoctorCount, request.SiteCount, request.AppCount, request.StartAt, request.EndAt, organizationID)
	var createVipQuery = "INSERT INTO `vip_schedule`(`doctor_count`,`site_count`,`app_count`, `start_at`, `end_at`, `organization_id`) VALUES (?,?,?,?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(createVipQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	res, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	lastID, _ := res.LastInsertId()
	v, _ := vip.GetVipScheduleByID(lastID)
	c.JSON(200, v)
}

func (oc *OrganizationControllerStruct) CreateOrganizationScheduleCase(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		return
	}
	var request organization.CreateVipScheduleCaseRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	v, _ := vip.GetVipScheduleCaseByName(request.Name)
	if v != nil {
		c.JSON(200, v)
		return
	}
	var values []interface{}
	values = append(values, request.Name, organizationID)
	var createVipQuery = "INSERT INTO `vip_cases`( `name`, `organization_id`) VALUES (?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(createVipQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	res, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	lastID, _ := res.LastInsertId()
	v, _ = vip.GetVipScheduleCaseByID(lastID)
	c.JSON(200, v)
}
