package easecrawler

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// 全局日志实例，默认输出到标准输出
var Logger *EaseLogger = NewLogger(os.Stdout)

type EaseLogger struct {
	logWriter *log.Logger
	mu        sync.Mutex // 写操作使用互斥锁，而非读写锁
	out       io.Writer
}

// NewLogger 创建日志实例，prefix 为全局前缀（如 [easecrawler]）
func NewLogger(out io.Writer) *EaseLogger {
	if out == nil {
		out = os.Stdout // 兜底，避免 nil writer
	}
	return &EaseLogger{
		logWriter: log.New(out, "[easecrawler] ", log.LstdFlags), // 全局前缀+时间戳
		mu:        sync.Mutex{},
		out:       out,
	}
}

// InitLogger 初始化全局日志实例
func InitLogger(out io.Writer) {
	Logger = NewLogger(out)
}

// 通用日志方法，提取重复逻辑
func (l *EaseLogger) log(level string, format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// 拼接级别和格式化字符串，级别后加空格提升可读性
	fullFormat := fmt.Sprintf("[%s] %s", level, format)
	l.logWriter.Printf(fullFormat, v...)
}

// ========== INFO 级别日志 ==========
func (l *EaseLogger) Printf(format string, v ...any) {
	l.log("INFO", format, v...)
}

func (l *EaseLogger) Println(v ...any) {
	l.log("INFO", "%v", v...)
}

func (l *EaseLogger) Print(v ...any) {
	l.log("INFO", "%v", v...)
}

// ========== ERROR 级别日志 ==========
func (l *EaseLogger) Errorf(format string, v ...any) {
	l.log("ERROR", format, v...)
}

func (l *EaseLogger) Errorln(v ...any) {
	l.log("ERROR", "%v", v...)
}

func (l *EaseLogger) Error(v ...any) {
	l.log("ERROR", "%v", v...)
}

// ========== WARN 级别日志 ==========
func (l *EaseLogger) Warnf(format string, v ...any) {
	l.log("WARN", format, v...)
}

func (l *EaseLogger) Warnln(v ...any) {
	l.log("WARN", "%v", v...)
}

func (l *EaseLogger) Warn(v ...any) {
	l.log("WARN", "%v", v...)
}
