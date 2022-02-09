package httpEngine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/appointmentController"
	authController "gitlab.com/simateb-project/simateb-backend/controller/authController"
	casesController "gitlab.com/simateb-project/simateb-backend/controller/caseController"
	"gitlab.com/simateb-project/simateb-backend/controller/casePriceController"
	"gitlab.com/simateb-project/simateb-backend/controller/caseTypeController"
	HoldayController "gitlab.com/simateb-project/simateb-backend/controller/holidayController"
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
	"gitlab.com/simateb-project/simateb-backend/controller/paymentController"
	"gitlab.com/simateb-project/simateb-backend/controller/smsController"
	"gitlab.com/simateb-project/simateb-backend/controller/transferController"
	"gitlab.com/simateb-project/simateb-backend/controller/uploadController"
	"gitlab.com/simateb-project/simateb-backend/controller/userController"
	userOrganizationPrice "gitlab.com/simateb-project/simateb-backend/controller/userOrganizationPriceController"
	"gitlab.com/simateb-project/simateb-backend/middleware"
	"gitlab.com/simateb-project/simateb-backend/repository/appointment"
)

func Run(Port string) {
	engine := gin.Default()

	engine.Use(gin.Recovery())

	engine.Static("/images", "./images")

	//engine.Use(cors.New(cors.Config{
	//	//AllowOrigins:     []string{"*"},
	//	AllowAllOrigins:  true,
	//	AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "PATCH"},
	//	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "X-Requested-With", "Authorization"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//}))
	engine.Use(middleware.CORSMiddleware)

	v1 := engine.Group("api/v1")

	ac := authController.NewAuthController()
	oc := organizationController.NewOrganizationController()
	appr := appointment.NewAppointmentRepository()
	appc := appointmentController.NewAppointmentController(appr)
	uc := userController.NewUserController()
	ctc := caseTypeController.NewCaseTypeController(appr)
	hc := HoldayController.NewHolidayController()
	smsc := smsController.NewSMSController()
	notifc := smsController.NewNotificationController()
	pc := paymentController.NewPaymentController()
	caseC := casesController.NewCasesController()
	cpc := casePriceController.NewCasePriceController()
	uopc := userOrganizationPrice.NewUserOrganizationPriceController()
	uploadC := uploadController.NewUploadController()
	tc := transferController.NewTransferController()

	{
		v1.POST("/auth/login", ac.Login)
		v1.POST("/auth/login/code", ac.LoginWithCode)
	}

	{
		v1.POST("/upload", middleware.GinJwtAuth(uploadC.Upload, true, false))
		v1.POST("/upload/multiple", middleware.GinJwtAuth(uploadC.UploadMultipleFile, true, false))
		v1.GET("/file/organizations/:id/:name", uploadC.GetUploadedOrgImage)
		v1.GET("/file/appointments/:id/:prof", appc.GetResultImages)
		v1.GET("/file/results/:id/:prof/:name", uploadC.GetUploadedResultImage)
		v1.GET("/file/:path/:name", uploadC.GetUploadedFile)
	}

	{
		v1.POST("/organizations", middleware.GinJwtAuth(oc.Create, true, false))
		v1.GET("/organizations", middleware.GinJwtAuth(oc.GetList, true, false))
		v1.POST("/organizations/items", middleware.GinJwtAuth(oc.DeleteItems, true, false))
		v1.GET("/organizations/:id/holidays", middleware.GinJwtAuth(hc.GetList, true, false))
		v1.GET("/organizations/:id/radios", middleware.GinJwtAuth(appc.GetRadiosAppointmentList, true, false))
		v1.GET("/organizations/:id/photos", middleware.GinJwtAuth(appc.GetPhotosAppointmentList, true, false))
		v1.GET("/organizations/:id/offs", middleware.GinJwtAuth(appc.GetOffsAppointmentList, true, false))
		v1.GET("/organizations/all", middleware.GinJwtAuth(oc.GetListAll, true, false))
		v1.GET("/organizations/:id", middleware.GinJwtAuth(oc.Get, true, false))
		v1.GET("/organizations/:id/work-time", middleware.GinJwtAuth(oc.GetOrganizationWorkTime, true, false))
		v1.POST("/organizations/:id/work-time", middleware.GinJwtAuth(oc.UpdateOrganizationWorkTime, true, false))
		v1.PUT("/organizations/:id", middleware.GinJwtAuth(oc.Update, true, false))
		v1.GET("/organizations/:id/appointments", middleware.GinJwtAuth(oc.GetOrganizationAppointments, true, false))
		v1.GET("/organizations/:id/transfers", middleware.GinJwtAuth(oc.GetOrganizationTransferList, true, false))
		v1.GET("/organizations/:id/reports", middleware.GinJwtAuth(oc.GetOrganizationReports, true, false))
		v1.GET("/organizations/:id/images", middleware.GinJwtAuth(oc.GetOrganizationImages, true, false))
		v1.POST("/organizations/:id/images", middleware.GinJwtAuth(oc.UploadOrganizationImage, true, false))
		v1.GET("/organizations/:id/about", middleware.GinJwtAuth(oc.GetOrganizationAbout, true, false))
		v1.PUT("/organizations/:id/about", middleware.GinJwtAuth(oc.UpdateOrganizationAbout, true, false))
		v1.GET("/organizations/:id/rels", middleware.GinJwtAuth(oc.GetOrganizationRelList, true, false))
		v1.GET("/organizations/:id/vip/schedule", middleware.GinJwtAuth(oc.GetOrganizationScheduleList, true, false))
		v1.POST("/organizations/:id/vip/schedule", middleware.GinJwtAuth(oc.CreateOrganizationSchedule, true, false))
		v1.GET("/organizations/:id/vip/cases", middleware.GinJwtAuth(oc.GetOrganizationScheduleCasesList, true, false))
		v1.POST("/organizations/:id/vip/cases", middleware.GinJwtAuth(oc.CreateOrganizationScheduleCase, true, false))
		v1.GET("/organizations/:id/rel-offices", middleware.GinJwtAuth(oc.GetOrganizationRelOfficesList, true, false))
		v1.POST("/organizations/:id/slider", middleware.GinJwtAuth(oc.SetOrganizationSlider, true, false))
		v1.GET("/organizations/:id/users", middleware.GinJwtAuth(oc.GetUsers, true, false))
		v1.GET("/organizations/:id/employees", middleware.GinJwtAuth(oc.GetEmployees, true, false))
		v1.GET("/organizations/:id/wallet", middleware.GinJwtAuth(oc.GetOrganizationWallet, true, false))
		v1.GET("/organizations/:id/wallet/histories", middleware.GinJwtAuth(oc.GetOrganizationWalletHistories, true, false))
		v1.GET("/organizations/:id/wallet/histories/sum", middleware.GinJwtAuth(oc.GetOrganizationWalletHistoriesSum, true, false))
		v1.GET("/organizations/:id/wallet/histories/days", middleware.GinJwtAuth(oc.GetOrganizationWalletHistoriesDays, true, false))
		v1.POST("/organizations/:id/wallet/increase", middleware.GinJwtAuth(oc.IncreaseOrganizationWallet, true, false))
		v1.POST("/organizations/:id/wallet/decrease", middleware.GinJwtAuth(oc.DecreaseOrganizationWallet, true, false))
		v1.POST("/organizations/:id/wallet/set", middleware.GinJwtAuth(oc.SetOrganizationWallet, true, false))
	}

	{
		//v1.GET("/vip/schedule/reserve", middleware.GinJwtAuth(oc.ReserveVip, true, false))
		v1.GET("/vip/schedule/:id", middleware.GinJwtAuth(oc.GetOrganizationSchedule, true, false))
		v1.GET("/vip/cases/:id", middleware.GinJwtAuth(oc.GetVipScheduleCase, true, false))
	}

	{
		v1.GET("/own", middleware.GinJwtAuth(uc.Own, true, false))
		v1.POST("/users", middleware.GinJwtAuth(uc.Create, true, false))
		v1.GET("/users", middleware.GinJwtAuth(uc.GetList, true, false))
		v1.GET("/users/last-login", middleware.GinJwtAuth(uc.GetLastLoginUsers, true, false))
		v1.GET("/users/last-patient-login", middleware.GinJwtAuth(uc.GetLastLoginPatients, true, false))
		v1.POST("/users/items", middleware.GinJwtAuth(uc.DeleteItems, true, false))
		v1.PUT("/users/:id/code", middleware.GinJwtAuth(uc.CreateCode, true, false))
		v1.GET("/users/:id/transfers", middleware.GinJwtAuth(tc.GetUserTransferList, true, false))
		v1.GET("/users/:id/docs", middleware.GinJwtAuth(uc.GetUserDocs, true, false))
		v1.POST("/users/:id/docs", middleware.GinJwtAuth(uc.SendDoc, true, false))
		v1.GET("/users/:id", middleware.GinJwtAuth(uc.Get, true, false))
		v1.PUT("/users/:id", middleware.GinJwtAuth(uc.Update, true, false))
		v1.DELETE("/users/:id", middleware.GinJwtAuth(uc.Delete, true, false))
		v1.PATCH("/users/:id/password", middleware.GinJwtAuth(uc.ChangePassword, true, false))
		v1.GET("/users/:id/appointments", middleware.GinJwtAuth(uc.GetUserAppointmentList, true, false))
		v1.GET("/users/:id/appointments/results", middleware.GinJwtAuth(uc.GetUserAppointmentResultList, true, false))
		v1.GET("/users/:id/wallet", middleware.GinJwtAuth(uc.GetUserWallet, true, false))
		v1.GET("/users/:id/wallet/histories", middleware.GinJwtAuth(uc.GetUserWalletHistories, true, false))
		v1.POST("/users/:id/wallet/histories", middleware.GinJwtAuth(uc.CreateUserWalletHistories, true, false))
		v1.GET("/users/:id/wallet/histories/sum", middleware.GinJwtAuth(uc.GetUserWalletHistoriesSum, true, false))
		v1.GET("/users/:id/wallet/histories/days", middleware.GinJwtAuth(uc.GetUserWalletHistoriesDays, true, false))
		v1.POST("/users/:id/wallet/increase", middleware.GinJwtAuth(uc.IncreaseUserWallet, true, false))
		v1.POST("/users/:id/wallet/decrease", middleware.GinJwtAuth(uc.DecreaseUserWallet, true, false))
		v1.POST("/users/:id/wallet/set", middleware.GinJwtAuth(uc.SetUserWallet, true, false))
		v1.GET("/users/:id/payments", middleware.GinJwtAuth(pc.GetPaymentList, true, false))
		v1.GET("/users/:id/orthodontics", middleware.GinJwtAuth(uc.GetOrthodonticsList, true, false))
		v1.POST("/users/:id/orthodontics", middleware.GinJwtAuth(uc.CreateOrthodonticsList, true, false))
		v1.PUT("/docs/:id", middleware.GinJwtAuth(uc.UpdateDoc, true, false))
		v1.DELETE("/docs/:id", middleware.GinJwtAuth(uc.DeleteDoc, true, false))
	}

	{
		v1.GET("/transfers/:id", middleware.GinJwtAuth(tc.Get, true, false))
	}

	{
		v1.POST("/payments", middleware.GinJwtAuth(pc.Create, true, false))
		v1.POST("/payments/delete", middleware.GinJwtAuth(pc.DeleteItems, true, false))
		v1.PUT("/payments/:id", middleware.GinJwtAuth(pc.Update, true, false))
		v1.DELETE("/payments/:id", middleware.GinJwtAuth(pc.Delete, true, false))
	}

	{
		v1.GET("/appointments", middleware.GinJwtAuth(appc.GetAppointmentList, true, false))
		v1.POST("/appointments", middleware.GinJwtAuth(appc.Create, true, false))
		v1.GET("/appointments/code", middleware.GinJwtAuth(appc.GetAppointmentByCode, true, false))
		v1.GET("/appointments/:id", middleware.GinJwtAuth(appc.Get, true, false))
		v1.PUT("/appointments/:id", middleware.GinJwtAuth(appc.Update, true, false))
		v1.PATCH("/appointments/:id", middleware.GinJwtAuth(appc.ChangeStatus, true, false))
		v1.POST("/appointments/:id/results", middleware.GinJwtAuth(appc.SendResult, true, false))
		v1.POST("/appointments/:id/accept", middleware.GinJwtAuth(appc.AcceptAppointment, true, false))
		v1.POST("/appointments/:id/cancel", middleware.GinJwtAuth(appc.CancelAppointment, true, false))
		v1.GET("/appointments/:id/admissions", middleware.GinJwtAuth(appc.GetAppointmentAdmissions, true, false))
		v1.POST("/appointments/:id/finished", middleware.GinJwtAuth(appc.FinishAppointmentAdmissions, true, false))
		v1.GET("/appointments/que", middleware.GinJwtAuth(appc.GetQueDetails, true, false))
		v1.GET("/appointments/search", middleware.GinJwtAuth(appc.SearchAppointment, true, false))
	}

	{
		v1.GET("/holidays", middleware.GinJwtAuth(hc.GetList, true, false))
		v1.POST("/holidays", middleware.GinJwtAuth(hc.Create, true, false))
		v1.GET("/holidays/:id", middleware.GinJwtAuth(hc.Get, true, false))
		v1.PUT("/holidays/:id", middleware.GinJwtAuth(hc.Update, true, false))
		v1.DELETE("/holidays/:id", middleware.GinJwtAuth(hc.Delete, true, false))
	}

	{
		v1.GET("/case-types", middleware.GinJwtAuth(ctc.GetListByOrganization, true, false))
		v1.POST("/case-types", middleware.GinJwtAuth(ctc.Create, true, false))
		v1.GET("/case-types/:id", middleware.GinJwtAuth(ctc.Get, true, false))
		v1.PUT("/case-types/:id", middleware.GinJwtAuth(ctc.Update, true, false))
		v1.DELETE("/case-types/:id", middleware.GinJwtAuth(ctc.Delete, true, false))
	}

	{
		v1.GET("/cases", middleware.GinJwtAuth(caseC.GetList, true, false))
		v1.GET("/cases/insert", middleware.GinJwtAuth(caseC.InsertAll, true, false))
		v1.GET("/cases/prof", middleware.GinJwtAuth(caseC.GetList, true, false))
		v1.GET("/cases/:id", middleware.GinJwtAuth(caseC.Get, true, false))
	}

	{
		v1.GET("/sms", middleware.GinJwtAuth(smsc.GetList, true, false))
		v1.POST("/sms", middleware.GinJwtAuth(smsc.Create, true, false))
		v1.GET("/sms/:id", middleware.GinJwtAuth(smsc.Get, true, false))
		v1.POST("/sms/items", middleware.GinJwtAuth(smsc.Delete, true, false))
	}

	{
		v1.GET("/notifications", middleware.GinJwtAuth(notifc.GetList, true, false))
		v1.POST("/notifications", middleware.GinJwtAuth(notifc.Create, true, false))
		v1.GET("/notifications/:id", middleware.GinJwtAuth(notifc.Get, true, false))
		v1.POST("/notifications/items", middleware.GinJwtAuth(notifc.Delete, true, false))
	}

	{
		v1.GET("/operations", middleware.GinJwtAuth(appc.GetOperationList, true, false))
	}

	{
		v1.GET("/admin/users", middleware.GinJwtAuth(uc.GetListForAdmin, true, false))
		v1.GET("/admin/admins", middleware.GinJwtAuth(uc.GetAdminList, true, false))
		v1.POST("/admin/users", middleware.GinJwtAuth(uc.Create, true, false))
		v1.GET("/admin/users/:id", middleware.GinJwtAuth(uc.Get, true, false))
		v1.PUT("/admin/users/:id", middleware.GinJwtAuth(uc.Update, true, false))
		v1.DELETE("/admin/users/:id", middleware.GinJwtAuth(uc.Delete, true, false))
		v1.PATCH("/admin/users/:id/password", middleware.GinJwtAuth(uc.ChangePassword, true, false))
		v1.POST("/admin/users/:id/wallet/increase", middleware.GinJwtAuth(uc.IncreaseUserWallet, true, false))
		v1.POST("/admin/users/:id/wallet/decrease", middleware.GinJwtAuth(uc.DecreaseUserWallet, true, false))
		v1.POST("/admin/users/:id/wallet/set", middleware.GinJwtAuth(uc.SetUserWallet, true, false))

		v1.GET("/admin/holidays", middleware.GinJwtAuth(hc.GetList, true, false))
		v1.POST("/admin/holidays", middleware.GinJwtAuth(hc.Create, true, false))
		v1.GET("/admin/holidays/:id", middleware.GinJwtAuth(hc.Get, true, false))
		v1.PUT("/admin/holidays/:id", middleware.GinJwtAuth(hc.Update, true, false))
		v1.DELETE("/admin/holidays/:id", middleware.GinJwtAuth(hc.Delete, true, false))

		v1.GET("/admin/sms", middleware.GinJwtAuth(smsc.GetListForAdmin, true, false))
		v1.POST("/admin/sms", middleware.GinJwtAuth(smsc.Create, true, false))
		v1.GET("/admin/sms/:id", middleware.GinJwtAuth(smsc.Get, true, false))
		v1.DELETE("/admin/sms", middleware.GinJwtAuth(smsc.Delete, true, false))

		v1.POST("/admin/organizations", middleware.GinJwtAuth(oc.Create, true, false))
		v1.GET("/admin/organizations", middleware.GinJwtAuth(oc.GetListForAdmin, true, false))
		v1.GET("/admin/organizations/:id", middleware.GinJwtAuth(oc.Get, true, false))
		v1.PUT("/admin/organizations/:id", middleware.GinJwtAuth(oc.Update, true, false))
		v1.GET("/admin/organizations/:id/users", middleware.GinJwtAuth(oc.GetUsers, true, false))
		v1.GET("/admin/organization/:id/user-prices", middleware.GinJwtAuth(uopc.GetList, true, false))
		v1.POST("/admin/user-prices", middleware.GinJwtAuth(uopc.Create, true, false))
		v1.GET("/admin/user-prices/:id", middleware.GinJwtAuth(uopc.Get, true, false))
		v1.PUT("/admin/user-prices/:id", middleware.GinJwtAuth(uopc.Update, true, false))

		v1.GET("/admin/case-prices", middleware.GinJwtAuth(cpc.GetList, true, false))
		v1.POST("/admin/case-prices", middleware.GinJwtAuth(cpc.Create, true, false))
		v1.GET("/admin/case-prices/:id", middleware.GinJwtAuth(cpc.Get, true, false))
		v1.PUT("/admin/case-prices/:id", middleware.GinJwtAuth(cpc.Update, true, false))

		v1.GET("/admin/cases", middleware.GinJwtAuth(caseC.GetList, true, false))
		v1.POST("/admin/cases", middleware.GinJwtAuth(caseC.Create, true, false))
		v1.GET("/admin/cases/:id", middleware.GinJwtAuth(caseC.Get, true, false))
		v1.PUT("/admin/cases/:id", middleware.GinJwtAuth(caseC.Update, true, false))

		v1.POST("/admin/histories/accept", middleware.GinJwtAuth(uc.AcceptOrRejectWalletHistories, true, false))
		v1.GET("/admin/histories", middleware.GinJwtAuth(uc.GetWalletHistoriesForAdmin, true, false))
	}

	fmt.Println(engine.Run(fmt.Sprintf(":%s", Port)))
}
