package easecrawler

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Engine struct {
	CrawlerGroup
	ctx          *Context
	mu           *sync.RWMutex
	globalLogger *log.Logger
}

type registeredCrawler struct {
	path    string
	name    string
	crawler Crawler
	log     *log.Logger
}

type CrawlerGroup struct {
	path     string
	engine   *Engine
	children map[string]*CrawlerGroup
	parent   *CrawlerGroup
	crawlers map[string]*registeredCrawler
}

func New() *Engine {
	e := &Engine{
		CrawlerGroup: CrawlerGroup{
			path:     "/",
			children: make(map[string]*CrawlerGroup),
			parent:   nil,
			crawlers: make(map[string]*registeredCrawler),
		},
		mu:           new(sync.RWMutex),
		ctx:          new(Context),
		globalLogger: log.New(os.Stdout, "", log.LstdFlags),
	}
	e.engine = e
	return e
}

func (g *CrawlerGroup) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	g.RunWithContext(context.Background())
}

func (g *CrawlerGroup) RunWithContext(ctx context.Context) {
	g.engine.mu.RLock()
	items := make([]*registeredCrawler, 0, len(g.engine.crawlers))
	for _, c := range g.engine.crawlers {
		items = append(items, c)
	}
	g.engine.mu.RUnlock()

	for _, item := range items {
		go g.engine.runCrawlerLoop(ctx, item)
	}
}

func (e *Engine) runCrawlerLoop(ctx context.Context, item *registeredCrawler) {
	meta := item.crawler.Meta()
	interval := meta.Interval
	if interval <= 0 {
		interval = time.Second
	}

	runOnce := func() {
		logger := item.log
		cctx := &Context{}
		cctx.Set(ContextLoggerKey, logger)
		err := item.crawler.Run(cctx)
		if err != nil {
			e.logf(item.path, item.name, "run error: %v", err)
			return
		}
		e.logf(item.path, item.name, "run ok")
	}

	if meta.StartImmediately {
		runOnce()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			e.logf(item.path, item.name, "stopped")
			return
		case <-ticker.C:
			runOnce()
		}
	}
}

func (g *CrawlerGroup) Register(crawler Crawler) {
	name := crawler.Name()
	path := filepath.Join(g.path, name)

	g.engine.mu.Lock()
	defer g.engine.mu.Unlock()
	f, err := os.OpenFile(filepath.Join("/logs", path+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	g.crawlers[name] = &registeredCrawler{path: path, name: name, crawler: crawler, log: log.New(f, "", log.LstdFlags)}
}

func (g *CrawlerGroup) Group(name string) *CrawlerGroup {
	newGroup := &CrawlerGroup{
		path:     filepath.Join(g.path, name),
		engine:   g.engine,
		parent:   g,
		children: make(map[string]*CrawlerGroup),
		crawlers: make(map[string]*registeredCrawler),
	}
	g.children[name] = newGroup
	err := os.MkdirAll(filepath.Join("/logs", newGroup.path), 0755)
	if err != nil {
		panic(err)
	}
	return newGroup
}

func (e *Engine) logf(path, name, format string, args ...any) {
	e.globalLogger.Printf("[%s][%s] %s", path, name, fmt.Sprintf(format, args...))
}
