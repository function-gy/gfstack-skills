// ================================================================
// 示例: Captcha 验证码 (internal/library/captcha/captcha.go)
// 支持算数验证码和字符验证码
// ================================================================

package captcha

import (
	"context"
	"image/color"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/mojocn/base64Captcha"

	"hotgo/internal/consts"
)

var store = base64Captcha.DefaultMemStore

// Generate 生成验证码
func Generate(ctx context.Context, captchaType int) (id string, base64 string) {
	var err error

	switch captchaType {
	case consts.CaptchaTypeMath: // 算数
		driver := &base64Captcha.DriverMath{
			Height:     42,
			Width:      100,
			NoiseCount: 0,
			BgColor:    &color.RGBA{R: 255, G: 250, B: 250, A: 250},
			Fonts:      []string{"chromohv.ttf"},
		}
		c := base64Captcha.NewCaptcha(driver.ConvertFonts(), store)
		id, base64, _, err = c.Generate()

	default: // 字符
		driver := &base64Captcha.DriverString{
			Height: 42,
			Width:  100,
			Length: 4,
			BgColor: &color.RGBA{R: 255, G: 250, B: 250, A: 250},
			Source: "abcdefghjkmnpqrstuvwxyz23456789",
			Fonts:  []string{"chromohv.ttf"},
		}
		c := base64Captcha.NewCaptcha(driver.ConvertFonts(), store)
		id, base64, _, err = c.Generate()
	}

	if err != nil {
		g.Log().Errorf(ctx, "captcha.Generate err:%+v", err)
	}
	return
}

// Verify 验证验证码
func Verify(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	return store.Verify(id, gstr.ToLower(answer), true)
}
