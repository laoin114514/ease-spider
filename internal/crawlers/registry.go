package crawlers

import (
	"context"
	"fmt"
	"log"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type task1 struct{}

func (t *task1) Name() string {
	return "task1"
}

func (t *task1) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Group:            "task1",
		Interval:         1 * time.Second,
		StartImmediately: true,
	}
}

func (t *task1) Run(ctx context.Context, c *easecrawler.Context) error {
	if logger, ok := easecrawler.GetAs[*log.Logger](c, easecrawler.ContextLoggerKey); ok {
		logger.Println("task1 running")
	} else {
		fmt.Println("task1 running")
	}
	return nil
}

type task2 struct{}

func (t *task2) Name() string {
	return "task2"
}

func (t *task2) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Group:            "task2",
		Interval:         5 * time.Second,
		StartImmediately: true,
	}
}

func (t *task2) Run(ctx context.Context, c *easecrawler.Context) error {
	if logger, ok := easecrawler.GetAs[*log.Logger](c, easecrawler.ContextLoggerKey); ok {
		logger.Println("task2 running")
	} else {
		fmt.Println("task2 running")
	}
	return nil
}

func Run() {
	engine := easecrawler.NewEngine()
	defer engine.Close()
	g1 := engine.Group("g1")
	g11 := g1.Group("g11")
	g11.Register(&task1{})
	g2 := g1.Group("g2")
	g2.Register(&task2{})
	engine.Run()
}
