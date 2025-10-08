package services

import (
	"os"
	"spider/src/utils"
)

// BaseService 提供所有服务的公共功能
type LogService struct {
	log     *utils.LogContainer
	logPath string
	errPath string
}

// NewBaseService 创建基础服务
func NewLogService(logPath, errPath string) *LogService {
	return &LogService{
		log:     utils.NewLogContainer(),
		logPath: logPath,
		errPath: errPath,
	}
}

// GetLog 获取日志
func (b *LogService) GetLog() []string {
	return b.log.GetLog()
}

// GetErr 获取错误日志
func (b *LogService) GetErr() []string {
	return b.log.GetErr()
}

// SaveLog 保存日志
func (b *LogService) SaveLog() error {
	return SaveLog(b.logPath, b.log.GetLog())
}

// SaveErr 保存错误日志
func (b *LogService) SaveErr() error {
	return SaveErr(b.errPath, b.log.GetErr())
}

// Clear 清理日志
func (b *LogService) Clear() error {
	b.log.ClearLog()
	b.log.ClearErr()
	return nil
}

// AddLog 添加日志
func (b *LogService) AddLog(log string) {
	b.log.AddLog(log)
}

// AddErr 添加错误日志
func (b *LogService) AddErr(err string) {
	b.log.AddErr(err)
}
func SaveLog(logPath string, log []string) error {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, log := range log {
		file.WriteString(log)
	}
	return nil
}
func SaveErr(errPath string, errs []string) error {
	file, err := os.OpenFile(errPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, err := range errs {
		file.WriteString(err)
	}
	return nil
}
