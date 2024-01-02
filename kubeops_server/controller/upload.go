package controller

import (
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubeops/model"
	"kubeops/service"
	"os"
	"path"
	"strconv"
	"time"
)

type YamlData struct {
}

var Upload upload

type upload struct{}

const yamlurl = "./static/yaml/"

// 上传yaml文件创建资源
func (u *upload) uploadYamlFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["yamlfile"]
	for _, file := range files {
		logger.Info(file.Filename)
		// 上传文件至指定目录
		extname := path.Ext(file.Filename)
		exifile := map[string]bool{
			".yaml": true,
		}
		if exifile[extname] {
			dir := yamlurl + time.Now().Format("2006-01-02")
			if !model.DirExists(dir) {
				err := os.Mkdir(dir, os.ModePerm)
				if err != nil {
					c.JSON(400, gin.H{
						"code": 400,
						"msg":  "上传失败",
					})
					return
				}
			}
			dst := path.Join(dir, file.Filename)
			err := c.SaveUploadedFile(file, dst)
			if err != nil {
				logger.Info(err.Error())
				c.JSON(400, gin.H{
					"code": 400,
					"msg":  "上传失败",
				})
				return
			}
			uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
			msg, err := service.Upload.UploadFile(dst, uuid)
			if err != nil {
				c.JSON(400, gin.H{
					"code": 400,
					"msg":  err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  msg,
			})
			return
		}
		//dir := yamlurl + time.Now().Format("2006-01-02")
		////判断目录是不是存在,如果存在则不创建
		//if !dirExists(dir) {
		//	err := os.Mkdir(dir, os.ModePerm)
		//	if err != nil {
		//		c.JSON(400, gin.H{
		//			"code": 400,
		//			"msg":  "上传失败",
		//		})
		//		return
		//	}
		//}
		//dst := path.Join(dir, file.Filename)
		//err := c.SaveUploadedFile(file, dst)
		//if err != nil {
		//	logger.Info(err.Error())
		//	c.JSON(400, gin.H{
		//		"code": 400,
		//		"msg":  "上传失败",
		//	})
		//	return
		//}
		//c.JSON(200, gin.H{
		//	"code": 200,
		//	"msg":  "上传成功",
		//})
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "上传文件格式不正确",
		})
	}
}

// 通过 yaml 内容创建资源
func (u *upload) createYaml(c *gin.Context) {
	//定义并接受yaml 内容
	params := new(struct {
		YamlContent string `json:"yamlContent"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info(err)
	}
	// 将json 格式转化为yaml
	yamlStr, err := yaml.JSONToYAML([]byte(params.YamlContent))
	if err != nil {
		logger.Info("yaml 格式无效")
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "yaml 格式无效",
		})
		return
	}
	// 创建资源
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	msg, err := service.Upload.CreateYaml(yamlStr, uuid)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  msg,
	})
}
