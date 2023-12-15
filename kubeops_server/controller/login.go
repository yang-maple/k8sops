package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/wonderivan/logger"
	"image/color"
	"kubeops/service"
	"net/http"
)

type login struct{}

var Login login

var store = base64Captcha.DefaultMemStore

// CaptchaImage  生成验证码,并返回图片url
func (l *login) CaptchaImage(c *gin.Context) {
	//定义图片格式为 string
	driverConfig := base64Captcha.DriverString{
		Height:          32,                                     //高度
		Width:           110,                                    //宽度
		NoiseCount:      0,                                      //数字干扰项，数字越大，
		ShowLineOptions: 0,                                      //竖线干扰项
		Length:          4,                                      //文本长度
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm", //文本来源
		BgColor: &color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 1,
		},
		Fonts: []string{"wqy-microhei.ttc"}, // 文本样式可选
		// "3Dumb.ttf",
		// "ApothecaryFont.ttf",
		// "Comismsh.ttf",
		// "DENNEthree-dee.ttf",
		// "DeborahFancyDress.ttf",
		// "Flim-Flam.ttf",
		// "RitaSmith.ttf",
		// "actionj.ttf",
		// "chromohv.ttf",
		// "wqy-microhei.ttc",
	}
	//生成图片
	captcha := base64Captcha.NewCaptcha(driverConfig.ConvertFonts(), store)
	//获取图片id,url
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      1,
		"imageUrl":  b64s,
		"captchaId": id,
		"msg":       "success",
	})
}

// VerifyInfo 验证登陆信息
func (l *login) VerifyInfo(c *gin.Context) {
	params := new(service.LoginInfo)
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
	}
	//判断验证码是否正确 不正确直接返回
	if !store.Verify(params.CaptchaId, params.VerifyValue, true) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "验证码错误",
			"token": "",
		})
		return
	}
	// 验证码验证成功后，判断用户名密码
	token, err := service.Login.VerifyUserInfo(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   err.Error(),
			"token": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      params.Username + "登陆成功",
		"username": params.Username,
		"token":    token,
	})
}
