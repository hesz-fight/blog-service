package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/internal/service"
	"github.com/go-programming-tour/blog-service/pkg/app"
	"github.com/go-programming-tour/blog-service/pkg/convert"
	"github.com/go-programming-tour/blog-service/pkg/errcode"
)

// 标签处理器
type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

// @Summary 获取单个标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [get]
func (t *TagHandler) Get(c *gin.Context) {
}

// @Summary 获取多个标签
// @Produce  json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t *TagHandler) List(c *gin.Context) {
	param := service.TagListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 未通过参数验证
		global.Logger.ErrorfT("app.BindAndValid fail. errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	// 创建服务实例
	svc := service.New(c.Request.Context())
	// 获取标签数量
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
	if err != nil {
		global.Logger.ErrorfT("svc.Count fail. errs: %v", err)
		// 响应错误
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}
	// 获取分页数据
	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.ErrorfT("svc.GetTagList fail. errs: %v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}

	// 相应分页数据
	response.ToResponseList(tags, totalRows)
	return
}

// @Summary 新增标签
// @Produce  json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string true "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t *TagHandler) Create(c *gin.Context) {
	param := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 未通过参数验证
		global.Logger.ErrorfT("app.BindAndValid fail. errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&param)
	if err != nil {
		global.Logger.ErrorfT("svc.CreateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})
}

// @Summary 更新标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t *TagHandler) Update(c *gin.Context) {
	idStr := convert.StrTo(c.Param("id"))
	param := service.UpdateTagRequest{ID: idStr.MustUInt32()}

	fmt.Println(param)

	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 未通过参数验证
		global.Logger.ErrorfT("app.BindAndValid fail. errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context()) // 请求上下文
	err := svc.UpdateTag(&param)
	if err != nil {
		global.Logger.ErrorfT("svc.UpdateTag fail. err =", err)
		return
	}

	response.ToResponse(gin.H{}) // 响应数据

	return
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t *TagHandler) Delete(c *gin.Context) {
	idStr := convert.StrTo(c.Param("id"))
	param := service.DeleteTagRequest{ID: idStr.MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 未通过参数验证
		global.Logger.ErrorfT("app.BindAndValid fail. errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.ErrorfT("svc.DeleteTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}
