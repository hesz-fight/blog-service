package global

import (
	"github.com/go-programming-tour/blog-service/pkg/logger"
	"github.com/go-programming-tour/blog-service/pkg/setting"
)

// 三个区段配置文件的全局变量
var (
	ServerSetting   *setting.ServerSetting
	AppSetting      *setting.AppSetting
	DatabaseSetting *setting.DatabaseSetting
	JWTSetting      *setting.JWTSettingS
	EmailSetting    *setting.EmailSettingS
)

// 日志配置对象
var Logger *logger.Logger
