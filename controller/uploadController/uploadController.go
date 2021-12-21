package uploadController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	minio2 "gitlab.com/simateb-project/simateb-backend/utils/minio"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadControllerInterface interface {
	UploadByMINIO(c *gin.Context)
	Upload(c *gin.Context)
	UploadMultipleFile(c *gin.Context)
	GetList(c *gin.Context)
	Get(c *gin.Context)
	GetUploadedFile(c *gin.Context)
	GetUploadedOrgImage(c *gin.Context)
	GetUploadedResultImage(c *gin.Context)
}

type UploadControllerStruct struct {
}

func NewUploadController() UploadControllerInterface {
	x := &UploadControllerStruct{
	}
	return x
}

func (uc *UploadControllerStruct) UploadByMINIO(c *gin.Context) {
	err := minio2.MinioClient.MakeBucket(c, minio2.BucketName, minio.MakeBucketOptions{Region: minio2.Location})
	if err != nil {
		exists, errBucketExists := minio2.MinioClient.BucketExists(c, minio2.BucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", minio2.BucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", minio2.BucketName)
	}

	_, header, err := c.Request.FormFile("image")
	filename := header.Filename
	log.Println(filename)
	// Upload the zip file
	objectName := helper.RandomString(22)
	contentType := "image/png"

	// Upload the zip file with FPutObject
	info, err := minio2.MinioClient.FPutObject(c, minio2.BucketName, objectName, os.TempDir()+"\\"+filename, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	c.JSON(http.StatusOK, true)
}

func (uc *UploadControllerStruct) Upload(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err.Error(), "read file")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	t := time.Now().UnixNano()
	fileName := fmt.Sprintf("%d%s", t, filepath.Ext(header.Filename))
	staff := auth.GetStaffUser(c)
	path := fmt.Sprintf("./images/%d", staff.UserID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Println(err.Error())
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
	}
	err = c.SaveUploadedFile(header, path + "/" + fileName)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"name": fileName,
		"id": staff.UserID,
		"path": fmt.Sprintf("/images/%d/%s", staff.UserID, fileName),
	})
}

func (uc *UploadControllerStruct) UploadMultipleFile(c *gin.Context) {
	//files, header, err := c.Request.FormFile("files")
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	type UploadedFileStruct struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	files := form.File["files[]"]
	fileNames := []UploadedFileStruct{}
	var f UploadedFileStruct
	for _, file := range files {
		t := time.Now().UnixNano()
		fileName := fmt.Sprintf("%d%s", t , filepath.Ext(file.Filename))
		staff := auth.GetStaffUser(c)
		path := fmt.Sprintf("./images/%d", staff.UserID)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, 0755)
			if err != nil {
				log.Println(err.Error())
				errorsHandler.GinErrorResponseHandler(c, err)
				return
			}
		}
		err := c.SaveUploadedFile(file, path + "/" + fileName)
		if err != nil {
			log.Println(err.Error())
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		} // Err Handling
		f.Name = fileName
		f.Path = fmt.Sprintf("/images/%d/%s", staff.UserID, fileName)
		fileNames = append(fileNames, f)
	}
	c.JSON(http.StatusAccepted, gin.H{
		"names": fileNames,
	})
}

func (uc *UploadControllerStruct) Get(c *gin.Context) {

	c.JSON(http.StatusOK, true)
}

func (uc *UploadControllerStruct) GetList(c *gin.Context) {

	c.JSON(http.StatusOK, true)
}

func (uc *UploadControllerStruct) GetUploadedFile(c *gin.Context) {
	path := c.Param("path")
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"path": path,
		"url": fmt.Sprintf("http://%s/images/%s/%s", c.Request.Host, path, name),
	})
}
func (uc *UploadControllerStruct) GetUploadedOrgImage(c *gin.Context) {
	path := c.Param("id")
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"path": path,
		"url": fmt.Sprintf("http://%s/images/organizations/%s/%s", c.Request.Host, path, name),
	})
}
func (uc *UploadControllerStruct) GetUploadedResultImage(c *gin.Context) {
	id := c.Param("id")
	prof := c.Param("prof")
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"id": id,
		"url": fmt.Sprintf("http://%s/images/results/%s/%s/%s", c.Request.Host, id, prof, name),
	})
}
