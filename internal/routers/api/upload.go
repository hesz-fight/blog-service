package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/internal/service"
	"github.com/go-programming-tour/blog-service/pkg/app"
	"github.com/go-programming-tour/blog-service/pkg/convert"
	"github.com/go-programming-tour/blog-service/pkg/errcode"
	"github.com/go-programming-tour/blog-service/pkg/upload"
)

type UploadHandler struct {
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (u *UploadHandler) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails((err.Error())))
		return
	}

	tmp := convert.StrTo(c.PostForm("type"))
	fileType := tmp.MustInt()

	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
