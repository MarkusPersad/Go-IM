package validates

import (
	"Go-IM/pkg/err"
	"Go-IM/pkg/giutils"
	"Go-IM/pkg/zaplog"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// var Validatec *XValidator
//
//	func init() {
//		Validatec = New()
//	}
type (
	XValidator struct {
		validator *validator.Validate
	}
)

func New() *XValidator {
	valid := validator.New()
	for key, value := range rules {
		e := valid.RegisterValidation(key, value)
		if e != nil {
			zaplog.Logger.Fatal("validator register validation failed", zap.Error(e))
		}
	}
	return &XValidator{
		validator: valid,
	}
}

func (x *XValidator) Validate(data interface{}) error {
	errs := x.validator.Struct(data)
	if errs != nil {
		return err.NewError(1005, giutils.ProcessError(data, errs))
	}
	return nil
}
