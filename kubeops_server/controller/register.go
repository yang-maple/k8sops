package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
)

var Register register

type register struct{}

func (r *register) RegisterUser(c *gin.Context) {
	var registerInfo service.RegisterInfo
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	fmt.Println(&registerInfo)
	err = service.Register.RegisterUser(&registerInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
	})
}

func (r *register) SendEmail(c *gin.Context) {
	params := new(struct {
		Email string `json:"email"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	err = service.Register.SendEmail(params.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
	})
}
