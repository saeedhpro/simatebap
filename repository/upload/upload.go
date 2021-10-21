package upload

import "mime/multipart"

type Upload struct {
	Image *multipart.FileHeader `form:"image"`
}