package httpEngine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/appointmentController"
	authController "gitlab.com/simateb-project/simateb-backend/controller/authController"
	"gitlab.com/simateb-project/simateb-backend/controller/caseTypeController"
	HoldayController "gitlab.com/simateb-project/simateb-backend/controller/holidayController"
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
	"gitlab.com/simateb-project/simateb-backend/controller/smsController"
	"gitlab.com/simateb-project/simateb-backend/controller/userController"
	"gitlab.com/simateb-project/simateb-backend/middleware"
	"gitlab.com/simateb-project/simateb-backend/repository/appointment"
)

func Run(Port string) {
	engine := gin.Default()

	engine.Use(gin.Recovery())

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

	{
		v1.POST("/auth/login", ac.Login)
	}

	{
		v1.POST("/organizations", middleware.GinJwtAuth(oc.Create, true, false))
		v1.GET("/organizations", middleware.GinJwtAuth(oc.GetList, true, false))
		v1.GET("/organizations/:id", middleware.GinJwtAuth(oc.Get, true, false))
		v1.PUT("/organizations/:id", middleware.GinJwtAuth(oc.Update, true, false))
		v1.GET("/organizations/:id/users", middleware.GinJwtAuth(oc.GetUsers, true, false))
	}

	{
		v1.POST("/users", middleware.GinJwtAuth(uc.Create, true, false))
		v1.GET("/users", middleware.GinJwtAuth(uc.GetList, true, false))
		v1.GET("/users/:id", middleware.GinJwtAuth(uc.Get, true, false))
		v1.PUT("/users/:id", middleware.GinJwtAuth(uc.Update, true, false))
		v1.DELETE("/users/:id", middleware.GinJwtAuth(uc.Delete, true, false))
		v1.PATCH("/users/:id/password", middleware.GinJwtAuth(uc.ChangePassword, true, false))
	}

	{
		v1.GET("/appointments", middleware.GinJwtAuth(appc.GetAppointmentList, true, false))
		v1.POST("/appointments", middleware.GinJwtAuth(appc.Create, true, false))
		v1.PUT("/appointments/:id", middleware.GinJwtAuth(appc.Update, true, false))
		v1.PATCH("/appointments/:id", middleware.GinJwtAuth(appc.ChangeStatus, true, false))
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
		v1.GET("/sms", middleware.GinJwtAuth(smsc.GetList, true, false))
		v1.POST("/sms", middleware.GinJwtAuth(smsc.Create, true, false))
		v1.GET("/sms/:id", middleware.GinJwtAuth(smsc.Get, true, false))
		v1.DELETE("/sms", middleware.GinJwtAuth(smsc.Delete, true, false))
	}

	{
		v1.GET("/operations", middleware.GinJwtAuth(appc.GetOperationList, true, false))
	}

	{
		v1.GET("/admin/users", middleware.GinJwtAuth(uc.GetListForAdmin, true, false))
		v1.POST("/admin/users", middleware.GinJwtAuth(uc.Create, true, false))
		v1.GET("/admin/users/:id", middleware.GinJwtAuth(uc.Get, true, false))
		v1.PUT("/admin/users/:id", middleware.GinJwtAuth(uc.Update, true, false))
		v1.DELETE("/admin/users/:id", middleware.GinJwtAuth(uc.Delete, true, false))
		v1.PATCH("/admin/users/:id/password", middleware.GinJwtAuth(uc.ChangePassword, true, false))
		v1.GET("/admin/users/:id/appointments", middleware.GinJwtAuth(uc.GetUserAppointmentList, true, false))

		v1.GET("/admin/holidays", middleware.GinJwtAuth(hc.GetListForAdmin, true, false))
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
	}

	fmt.Println(engine.Run(fmt.Sprintf(":%s", Port)))
}
