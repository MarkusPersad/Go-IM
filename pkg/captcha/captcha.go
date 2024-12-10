package captcha

import (
	"Go-IM/pkg/common/customtypes"
	"Go-IM/pkg/err"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type Capt struct {
	store  base64Captcha.Store
	driver base64Captcha.Driver
}

func New(store base64Captcha.Store) Capt {
	mathDriver := base64Captcha.NewDriverMath(40, 160, 5, base64Captcha.OptionShowSineLine, &color.RGBA{
		R: 254,
		G: 254,
		B: 254,
		A: 254,
	}, base64Captcha.DefaultEmbeddedFonts, []string{"wqy-microhei.ttc"})
	return Capt{
		store:  store,
		driver: mathDriver,
	}
}

func (cap *Capt) Generate() (customtypes.CaptDataBase64, error) {
	c := base64Captcha.NewCaptcha(cap.driver, cap.store)
	id, b64s, _, e := c.Generate()
	if e != nil {
		return customtypes.CaptDataBase64{}, err.CheckCode
	}
	return customtypes.CaptDataBase64{
		Id:   id,
		B64s: b64s,
	}, nil
}

func (cap *Capt) Verify(id, answer string, clear bool) error {
	if len(answer) == 0 || len(id) == 0 {
		return err.CheckCode
	}
	c := base64Captcha.NewCaptcha(cap.driver, cap.store)
	match := c.Verify(id, answer, clear)
	if !match {
		return err.CheckCode
	}
	return nil
}
