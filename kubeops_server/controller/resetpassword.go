package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubeops/service"
	"net/http"
)

var ResetPassword resetPassword

type resetPassword struct{}

// verifyIdentity 身份验证
func (rs *resetPassword) verifyIdentity(c *gin.Context) {
	params := new(struct {
		Email      string `json:"email"`
		VerifyCode string `json:"verify_code"`
	})
	if err := c.ShouldBindJSON(&params); err != nil {
		logger.Info(errors.New("bind json error" + err.Error()))
		return
	}
	if err := service.ResetPassword.VerifyIdentity(params.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "邮箱验证失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码已发送至邮箱",
	})
}

// FindPass 找回并重置密码
func (rs *resetPassword) FindPass(c *gin.Context) {
	params := new(struct {
		Email      string `json:"email"`
		VerifyCode string `json:"verify_code"`
		Password   string `json:"password"`
	})
	if err := c.ShouldBindJSON(&params); err != nil {
		logger.Info(errors.New("bind json error" + err.Error()))
		return
	}
	if err := service.ResetPassword.FindPasswd(params.Email, params.VerifyCode, params.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "重置失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码重置成功",
	})

}

// ResetPass 重置密码
func (rs *resetPassword) ResetPass(c *gin.Context) {
	params := new(struct {
		Id        uint   `json:"id"`
		OldPasswd string `json:"old_passwd"`
		NewPasswd string `json:"new_passwd"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Error("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "绑定参数失败",
		})
		return
	}
	if err := service.ResetPassword.ResetPasswd(params.Id, params.OldPasswd, params.NewPasswd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码重置成功",
	})
}
