package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
)

// 校验错误结构体
type ValidError struct {
	Key     string
	Message string
}

func (v *ValidError) Error() string {
	return v.Message
}

// 校验错误结构体切片
type ValidErrors []*ValidError

func (v *ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v *ValidErrors) Errors() []string {
	var errs []string
	for _, err := range *v {
		errs = append(errs, err.Error())
	}

	return errs
}

// 绑定参数并且校验,
// 返回是否成功过通过验证的标记和验证错误.
func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v) // 绑定引擎
	if err != nil {
		v := c.Value("trans") //

		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			// 非验证错误直接返回
			return false, errs
		}

		// 翻译参数验证错误
		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{Key: key, Message: value})
		}

		return false, errs
	}

	return true, nil
}
