package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/internal/model"
	"github.com/go-programming-tour/blog-service/internal/routers"
	"github.com/go-programming-tour/blog-service/pkg/logger"
	"github.com/go-programming-tour/blog-service/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting fail. err = %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine fail. err = %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger fail. err = %v", err)
	}
}

// @title 博客系统
// @version 1.0
// @description Go 语言编程之旅：一起用 Go 做项目
// @termsOfService https://github.com/go-programming-tour-book
func main() {

	// global.Logger.InfofT("%s: go-programing-book/%s", "Testing", "blog_service")

	gin.SetMode(global.ServerSetting.RunMode)

	router := routers.NewRouter()

	server := &http.Server{
		Addr:           "127.0.0.1" + ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}

// 读取配置文件，返回配置参数结构体
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	// fmt.Println(*global.ServerSetting)
	// fmt.Println(*global.AppSetting)
	// fmt.Println(*global.DatabaseSetting)

	return nil
}

// 创建数据库实例
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngin(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

// 初始化日志组件
func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" +
		global.AppSetting.LogFileName + global.AppSetting.LogFileExt

	// 使用lumberjack作为日志库的io.Writer
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   600,
		MaxAge:    60,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}
