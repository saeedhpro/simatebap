package transferController

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/repository/transfer"
	"net/http"
	"strconv"
)

type TransferControllerInterface interface {
	Get(c *gin.Context)
	GetUserTransferList(c *gin.Context)
}

type TransferRepositoryStruct struct {
}

type TransferControllerStruct struct {
}

func NewTransferController() TransferControllerInterface {
	x := &TransferControllerStruct{
	}
	return x
}
func (uc *TransferControllerStruct) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	tid, _ := strconv.ParseInt(id, 10, 64)
	data, _ := transfer.GetTransferByID(tid)
	c.JSON(http.StatusOK, data)
}

func (uc *TransferControllerStruct) GetUserTransferList(c *gin.Context) {
	id := c.Param("id")
	page := c.Query("page")
	if id == "" {
		return
	}
	if page == "" {
		page = "1"
	}
	tid, _ := strconv.ParseInt(id, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)
	data := transfer.GetUserTransfers(tid, p)
	c.JSON(http.StatusOK, data)
}
