package easecrawler

import (
	"log"
	"sync"
	"time"
)

const ContextLoggerKey = "logger"

type Crawler interface {
	Run(c *Context) error
	Name() string
	Meta() Meta
}

type Meta struct {
	Interval         time.Duration
	StartImmediately bool
	//自定义日志(默认使用系统日志)
	Logger *log.Logger
}

type Context struct {
	values sync.Map
}

func (c *Context) Set(key string, value any) {
	c.values.Store(key, value)
}

func (c *Context) Get(key string) (any, bool) {
	return c.values.Load(key)
}

func GetAs[T any](c *Context, key string) (T, bool) {
	if c == nil {
		var zero T
		return zero, false
	}
	v, ok := c.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	val, ok := v.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return val, true
}
