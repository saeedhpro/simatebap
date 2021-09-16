package main

import (
	httpEngine "gitlab.com/simateb-project/simateb-backend/controller/http"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
)

func main() {
	repository.Init()
	errorsHandler.Init()
	httpEngine.Run("8000")
}
