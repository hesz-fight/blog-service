package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour/blog-service/pkg/errcode"
)

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `josn:"page_size"`
	TotalRows int `json:"total_rows"`
}

// 响应处理
type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"page": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{
		"Code": err.GetCode(),
		"msg":  err.GetMsg(),
	}
	details := err.GetDetails()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}
