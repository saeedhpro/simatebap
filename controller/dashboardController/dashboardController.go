package DashboardController

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/dashboard"
	"gitlab.com/simateb-project/simateb-backend/repository/dashboardRepository"
	"net/http"
)

type DashboardControllerInterface interface {
	GetMainData(c *gin.Context)
}

type DashboardControllerStruct struct {
}

func NewDashboardController() DashboardControllerInterface {
	x := &DashboardControllerStruct{
	}
	return x
}

func (uc *DashboardControllerStruct) GetMainData(c *gin.Context) {
	var dashboardData dashboard.DashboardData
	dashboardData.Count = dashboardRepository.GetTodayAppointments()
	lastLoginUsers, err := dashboardRepository.GetLastOnLineUsers()
	if err == nil {
		dashboardData.LastLoginUsers = lastLoginUsers
	}
	dashboardData.UnknownGenderUserCount = dashboardRepository.GetUnknownGenderUsersCount()
	dashboardData.FemaleGenderUserCount = dashboardRepository.GetFemaleUsersCount()
	dashboardData.MaleGenderUserCount = dashboardRepository.GetMaleUsersCount()

	c.JSON(http.StatusOK, dashboardData)
}
