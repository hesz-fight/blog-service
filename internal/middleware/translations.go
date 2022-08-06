package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 错误提示多语言的功能支持
func Translations() gin.HandlerFunc {

	return func(c *gin.Context) {
		uti := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		// 从请求头获取语言
		locale := c.GetHeader("locale")
		// 翻译器
		trans, _ := uti.GetTranslator(locale)
		// 验证器
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				// 注册
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			}

			c.Set("trans", trans) // 设置翻译器，key ---> value
		}

		c.Next()
	}
}
