package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
)

var Logger *logrus.Logger

type MyFormatter struct {
}

// Format 实现 Formatter 接口 Format方法返回自定义日志格式
func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	caller := getCallerInfo()
	return []byte(fmt.Sprintf("%s [%s] [%s:%d] %s\n", entry.Time.Format("2006-01-02 15:04:05"), strings.ToUpper(entry.Level.String()), caller.File, caller.Line, entry.Message)), nil
}

// InitLog 初始化日志格式
func InitLog() {
	//定义一个 loggers 实例
	Logger = logrus.New()
	// 设置日志级别
	Logger.SetLevel(logrus.DebugLevel)
	// 设置日志格式
	Logger.SetFormatter(new(MyFormatter))
	// 设置输出
	Logger.SetOutput(os.Stdout)
}

// getCallerInfo 获取调用者信息
func getCallerInfo() *runtime.Frame {
	pc, file, line, ok := runtime.Caller(7) // 通过runtime.Caller获取调用者信息
	if !ok {
		return nil
	}
	// 获取调用者信息
	caller := runtime.FuncForPC(pc)
	// 获取文件名和行号
	return &runtime.Frame{Function: path.Base(caller.Name()), File: file, Line: line}
}
