package utils

import (
	"fmt"
	"log"
	"time"
)

// 计时器
type Timer struct{}

func NewTimer() *Timer {
	return &Timer{}
}
func (t *Timer) CountDurationInMs(f func()) int64 {
	start := time.Now()
	f()
	end := time.Since(start).Milliseconds()
	return end
}
func (t *Timer) CountDurationInS(f func()) int64 {
	start := time.Now()
	f()
	end := int64(time.Since(start).Seconds())
	return end
}
func (t *Timer) CountDurationStr(f func()) string {
	start := time.Now()
	f()
	end := time.Since(start).String()
	return end
}
func (t *Timer) RunWithTimer(duration time.Duration, f func() error) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			err := f()
			if err != nil {
				log.Println(err)
			}
		}
	}()
	log.Println("定时器启动(单任务)")
}
func (t *Timer) MultiRunWithTimer(duration time.Duration, f ...func() error) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			for _, f := range f {
				go func() {
					err := f()
					if err != nil {
						log.Println(err)
					}
				}()
			}
		}
	}()
	log.Println("定时器启动(多任务)")
}

// 日期格式化工具
type DateFormat struct {
}

func NewDateFormat() *DateFormat {
	return &DateFormat{}
}
func (d *DateFormat) DateTimeWithSecond(rawStamp int64) string {
	stamp := time.Unix(rawStamp, 1)
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", stamp.Year(), stamp.Month(), stamp.Day(), stamp.Hour(), stamp.Minute(), stamp.Second())
	return str
}
func (d *DateFormat) NowDateTime() string {
	now := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return str
}
func (d *DateFormat) NowDate() string {
	now := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	return str
}
func (d *DateFormat) BeforDateTimeWithDay(day int64) string {
	now := time.Now()
	befor := time.Unix(now.Unix()-day*24*3600, 1)
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", befor.Year(), befor.Month(), befor.Day(), befor.Hour(), befor.Minute(), befor.Second())
	return str
}
