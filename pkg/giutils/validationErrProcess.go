package giutils

import (
	"Go-IM/pkg/common/defines"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func ProcessError(u interface{}, e error) string {
	if e == nil {
		return ""
	}
	var invalid *validator.InvalidValidationError
	if errors.As(e, &invalid) {
		return fmt.Sprintf("输入参数错误: %s", invalid.Error())
	}
	var validationErrs validator.ValidationErrors
	if errors.As(e, &validationErrs) {
		for _, validationError := range validationErrs {
			fieldName := validationError.Field()
			typeOf := reflect.TypeOf(u)
			if typeOf.Kind() == reflect.Pointer {
				typeOf = typeOf.Elem()
			}
			if field, o := typeOf.FieldByName(fieldName); o {
				errorInfo := field.Tag.Get(defines.FIELD_ERROR_INFO)
				return fmt.Sprintf("%s : %s", fieldName, errorInfo)
			} else {
				return "缺失字段错误信息"
			}
		}
	}
	return ""
}
