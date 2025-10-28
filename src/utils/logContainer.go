package utils

import "sync"

//============================日志容器===============================================
type LogContainer struct {
	logs       []string
	errs       []string
	mutex      sync.RWMutex
	dateFormat *DateFormat
}

func NewLogContainer() *LogContainer {
	return &LogContainer{
		logs:       []string{},
		errs:       []string{},
		dateFormat: NewDateFormat(),
	}
}
func (l *LogContainer) AddLog(log string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	log = l.dateFormat.NowDateTime() + " " + log
	l.logs = append(l.logs, log+"\n")
}
func (l *LogContainer) AddErr(err string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	err = l.dateFormat.NowDateTime() + " " + err
	l.errs = append(l.errs, err+"\n")
}
func (l *LogContainer) GetLog() []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.logs
}
func (l *LogContainer) GetErr() []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.errs
}
func (l *LogContainer) ClearLog() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logs = []string{}
}
func (l *LogContainer) ClearErr() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.errs = []string{}
}
