package mysqlQuery

const (
	LoginQuery                           = "SELECT id, fname, lname, tel, user_group_id, organization_id FROM user WHERE tel = ?"
	CreateOrganizationQuery              = "INSERT INTO `organization`(`name`, `known_as`, `profession_id`, `logo`,`phone`, `phone1`, `staff_id`, `info`, `case_types`, `sms_price`, `sms_credit`, `website`, `instagram`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"
	GetSimpleOrganizationQuery           = "SELECT id, ifnull(name, ''), profession_id FROM organization WHERE id = ?"
	GetAdminUsersQuery                   = "SELECT id, ifnull(fname, ''), ifnull(lname, ''), ifnull(tel, ''), ifnull(user_group_id, ''), ifnull(created, ''), ifnull(last_login, ''), ifnull(birth_date, '') FROM user LIMIT 10 OFFSET ?"
	GetUserGroupQuery                    = "SELECT id, ifnull(name, '') FROM user_group WHERE id = ?"
	GetSimpleProfessionQuery             = "SELECT id, ifnull(name, '') FROM profession WHERE id = ?"
	GetSimpleStaffQuery                  = "SELECT user.id as id, ifnull(user.fname, '') as fname, ifnull(user.lname, '') as lname, ifnull(organization.name, '') as organization FROM user LEFT JOIN organization ON user.organization_id = organization.id WHERE user.id = ?"
	CreateUserQuery                      = "INSERT INTO `user` (`fname`, `lname`, `email`, `info`, `description`, `file_id`, `gender`, `staff_id`, `user_group_id`, `organization_id`, `tel`, `tel1`, `nid`, `birth_date`, `address`, `introducer`, `pass`, `relation`, `logo`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	GetUserQuery                         = "SELECT id, organization_id, fname, lname, tel, user_group_id FROM `user` WHERE id = ?"
	GetUserOrganizationQuery             = "SELECT user.id id, ifnull(user.appcode, '') appcode,ifnull(user.logo, '') logo, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, user.last_login last_login, user.created created, user.tel tel, user.user_group_id user_group_id, user_group.name user_group_name, user.birth_date birth_date, organization.id organization_id, organization.name organization_name, user.relation relation, ifnull(user.description, '') description, ifnull(user.info, '') info, ifnull(user.tel1, '') tel1, ifnull(user.nid, '') nid, ifnull(user.address, '') address, ifnull(user.introducer, '') introducer, ifnull(user.gender, '') gender, ifnull(user.file_id, '') file_id FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id WHERE user.id = ?"
	DeleteUserQuery                      = "DELETE FROM `user` WHERE id = ?"
	ChangePasswordQuery                  = "UPDATE `user` SET `pass` = ? WHERE id = ?"
	GetOrganizationRelations             = "SELECT rel_organization_id FROM rel_organization where organization_id = ?"
	GetAppointmentListQuery              = "SELECT appointment.id id, appointment.case_type case_type, appointment.is_vip is_vip, appointment.start_at start_at, appointment.user_id user_id, appointment.info info, appointment.income income, appointment.status appointment_status,  appointment.updated_at updated_at, user.fname user_fname, user.lname user_lname, user.id user_id, user.gender user_gender from appointment LEFT JOIN user on appointment.user_id = user.id WHERE organization_id = ? and start_date >= ? and start_date < ?"
	GetAppointmentQuery                  = "SELECT appointment.id id, appointment.case_type case_type, appointment.is_vip is_vip, appointment.start_at start_at, appointment.user_id user_id, appointment.info info, appointment.income income, appointment.status appointment_status,  appointment.updated_at updated_at, user.fname user_fname, user.lname user_lname, user.id user_id, user.gender user_gender, appointment.price price, from appointment LEFT JOIN user on appointment.user_id = user.id WHERE id = ?"
	GetCaseTypesListForOrganizationQuery = "SELECT id, name, organization_id, duration, is_limited, limitation FROM case_type WHERE organization_id = ?"
	GetOrganizationOperationListQuery    = "SELECT id, user_id, start_at, info, income, case_type FROM `appointment` WHERE `case_type` = 'جراحی' and organization_id = ? AND start_at between ? AND ? ORDER BY `case_type` DESC"
	GetOrganizationAppointmentListQuery  = "SELECT id, user_id, start_at, info, income, case_type FROM `appointment` WHERE status in (?) and organization_id = ? AND start_at between ? AND ? ORDER BY `id` DESC"
	GetHolidayQuery                      = "SELECT holiday.id id, holiday.title, holiday.hdate, holiday.organization_id, organization.name organization_name FROM holiday LEFT JOIN organization ON holiday.organization_id = organization.id WHERE holiday.id = ?"
	CreateHolidayQuery                   = "INSERT INTO `holiday`(`hdate`, `organization_id`, `title`) VALUES (?,?,?)"
	DeleteHolidayQuery                   = "DELETE FROM `holiday` WHERE `id` = ?"
	CreateCaseTypeQuery                  = "INSERT INTO `case_type`(`name`, `organization_id`, `duration`, `is_limited`, `limitation`) VALUES (?,?,?,?,?)"
	GetCaseTypeQuery                     = "SELECT `id`, `name`, `organization_id`, `duration`, `is_limited`, `limitation` FROM `case_type` WHERE id = ?"
	DeleteCaseTypeQuery                  = "DELETE FROM `case_type` WHERE id = ?"
	CreateSMSQuery                       = "INSERT INTO `sms`(`user_id`, `staff_id`, `number`, `msg`, `sent`, `incoming`, `organization_id`) VALUES (?,?,?,?,?,?,?)"
)

// select user.id id, user.fname fname, user.lname lname, organization.id organization_id, organization.profession_id profession_id, user.user_group_id user_group_id, user_group.name user_group_name FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id
