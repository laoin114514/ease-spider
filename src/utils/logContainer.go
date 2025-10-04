package utils

import "sync"

type LogContainer struct {
	logs  []string
	errs  []string
	mutex sync.RWMutex
}

func NewLogContainer() *LogContainer {
	return &LogContainer{
		logs: []string{},
		errs: []string{},
	}
}
func (l *LogContainer) AddLog(log string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logs = append(l.logs, log+"\n")
}
func (l *LogContainer) AddErr(err string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
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
