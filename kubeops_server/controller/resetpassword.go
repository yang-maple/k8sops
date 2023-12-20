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
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "验证码已发送至邮箱",
	})
}

func (rs *resetPassword) reset(c *gin.Context) {
	params := new(struct {
		Email      string `json:"email"`
		VerifyCode string `json:"verify_code"`
		Password   string `json:"password"`
	})
	if err := c.ShouldBindJSON(&params); err != nil {
		logger.Info(errors.New("bind json error" + err.Error()))
		return
	}
	if err := service.ResetPassword.Reset(params.Email, params.VerifyCode, params.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "密码重置成功",
	})

}
