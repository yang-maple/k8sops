package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
)

var Register register

type register struct{}

// RegisterUser 注册用户
func (r *register) RegisterUser(c *gin.Context) {
	var registerInfo service.RegisterInfo
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	err = service.Register.RegisterUser(&registerInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "注册失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户" + registerInfo.Username + "注册成功",
	})
}

// SendEmail 注册用户时发送邮件
func (r *register) SendEmail(c *gin.Context) {
	params := new(struct {
		Email string `json:"email"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	err = service.Register.SendEmail(params.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "发送失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码已发送至邮箱",
	})
}
