package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"log"
	"os"
	"path"
	"time"
)

var Upload upload

type upload struct{}

const imageurl = "./static/image/"
const fileurl = "./static/file/"
const yamlurl = "./static/yaml/"

func (u *upload) uploadFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	logger.Info(files)
	for _, file := range files {
		logger.Info(file.Filename)
		log.Println(file.Filename)
		// 上传文件至指定目录
		extname := path.Ext(file.Filename)
		exifile := map[string]bool{
			".jpg": true,
			".png": true,
		}
		if exifile[extname] {
			dir := imageurl + time.Now().Format("2006-01-02")
			if !dirExists(dir) {
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
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "上传成功",
			})
		}
		dir := yamlurl + time.Now().Format("2006-01-02")
		//判断目录是不是存在,如果存在则不创建
		if !dirExists(dir) {
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
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功",
		})

	}
}

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsExist(err) {
		return false
	}
	return true
}
