package main

import (
	httpEngine "gitlab.com/simateb-project/simateb-backend/controller/http"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"gitlab.com/simateb-project/simateb-backend/utils/minio"
)

func main() {
	minio.Init()
	repository.Init()
	errorsHandler.Init()
	httpEngine.Run("8000")
}
