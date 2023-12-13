package controller

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubeops/service"
	"net/http"
	"time"
)

type login struct{}

var Login login

// CaptchaId png
var CaptchaId string

// CaptchaImage  根据id 获取验证码图片
func (l *login) CaptchaImage(c *gin.Context) {
	_ = serverHTTP(c.Writer, c.Request, CaptchaId, ".png", "zh", false, 120, 32)
}

// ImageId 获取验证码图片id
func (l *login) ImageId(c *gin.Context) {
	CaptchaId = captcha.NewLen(4)
	c.JSON(http.StatusOK, gin.H{
		"captchaId": CaptchaId,
	})
}

// VerifyInfo 验证登陆信息
func (l *login) VerifyInfo(c *gin.Context) {
	param := new(service.LoginInfo)
	err := c.ShouldBindJSON(&param)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
	}
	if !captcha.VerifyString(param.CaptchaId, param.CaptchaSolution) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "验证码错误",
			"token": "",
		})
	} else {
		token, err := service.Login.VerifyUserInfo(param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"token": nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"token": token,
		})
	}
}

// 重写captcha.server 根据捕获的标识返回图片或者音频
func serverHTTP(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
