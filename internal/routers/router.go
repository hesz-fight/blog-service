package routers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-programming-tour/blog-service/docs"
	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/internal/middleware"
	"github.com/go-programming-tour/blog-service/internal/routers/api"
	v1 "github.com/go-programming-tour/blog-service/internal/routers/api/v1"
	"github.com/go-programming-tour/blog-service/pkg/limiter"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	engin := gin.New()

	if global.ServerSetting.RunMode == "dubug" {
		// 注册日志中间件
		engin.Use(gin.Logger())
		// 注册异常捕捉中间件
		engin.Use(gin.Recovery())
	} else {
		engin.Use(middleware.AccessLog())
		engin.Use(middleware.Recovery())
	}

	engin.Use(middleware.RateLimiter(methodLimiters))
	engin.Use(middleware.ContextTimeout(60 * time.Second))

	// 注册翻译中间件
	engin.Use(middleware.Translations())
	// 注册swag路由
	engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册上传文件的路由
	upload := api.NewUploadHandler()
	engin.POST("/upload/file", upload.UploadFile)
	engin.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	// 注册jwt路由
	engin.POST("/auth", api.GetAuth)

	article := v1.NewArticleHandler()
	tag := v1.NewTagHandler()
	apiv1 := engin.Group("/api/v1")
	// 添加jwt验证
	apiv1.Use(middleware.JWT())
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles", article.Get)
		apiv1.GET("/articles/:id", article.List)
	}

	return engin
}
