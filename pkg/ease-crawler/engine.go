package easecrawler

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	CrawlerGroup
	ctx          *Context
	mu           sync.RWMutex
	globalLogger *log.Logger
	logRootDir   string
	crawlers     map[string]*registeredCrawler
	logFiles     map[string]*os.File
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
}

func New() *Engine {
	e := &Engine{
		CrawlerGroup: CrawlerGroup{
			path:     "/",
			children: make(map[string]*CrawlerGroup),
			parent:   nil,
		},
		mu:           sync.RWMutex{},
		ctx:          new(Context),
		globalLogger: log.New(os.Stdout, "", log.LstdFlags),
		logRootDir:   "logs",
		crawlers:     make(map[string]*registeredCrawler),
		logFiles:     make(map[string]*os.File),
	}
	e.engine = e
	return e
}

func (e *Engine) SetLogRootDir(dir string) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = "logs"
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.logRootDir = dir
}

func (g *CrawlerGroup) Run() {
	g.RunWithContext(context.Background())
}

func (g *CrawlerGroup) RunWithContext(ctx context.Context) {
	g.engine.mu.RLock()
	items := make([]*registeredCrawler, 0, len(g.engine.crawlers))
	for _, c := range g.engine.crawlers {
		items = append(items, c)
	}
	g.engine.mu.RUnlock()

	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(it *registeredCrawler) {
			defer wg.Done()
			g.engine.runCrawlerLoop(ctx, it)
		}(item)
	}
	wg.Wait()
}

func (e *Engine) runCrawlerLoop(ctx context.Context, item *registeredCrawler) {
	meta := item.crawler.Meta()
	interval := meta.Interval
	if interval <= 0 {
		interval = time.Second
	}

	runOnce := func() {
		logger := item.log
		if meta.Logger != nil {
			logger = meta.Logger
		}
		cctx := &Context{}
		cctx.Set(ContextLoggerKey, logger)
		err := item.crawler.Run(cctx)
		if err != nil {
			e.logf(item.path, item.name, "运行失败: %v", err)
			return
		}
		e.logf(item.path, item.name, "运行成功")
	}

	if meta.StartImmediately {
		runOnce()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			e.logf(item.path, item.name, "停止运行")
			return
		case <-ticker.C:
			runOnce()
		}
	}
}

func (g *CrawlerGroup) Register(crawler Crawler) {
	name := strings.TrimSpace(crawler.Name())
	if name == "" {
		g.engine.logf(g.path, "unknown", "注册crawler失败: 空名称")
		return
	}

	groupPath := cleanGroupPath(g.path)
	key := filepath.ToSlash(filepath.Join(groupPath, name))
	filePath := g.engine.crawlerLogPath(groupPath, name)

	g.engine.mu.Lock()
	defer g.engine.mu.Unlock()

	if _, exists := g.engine.crawlers[key]; exists {
		g.engine.logf(groupPath, name, "注册失败: 重复的crawler")
		return
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		g.engine.logf(groupPath, name, "注册失败: 创建目录失败: %v", err)
		return
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		g.engine.logf(groupPath, name, "注册失败: 打开日志文件失败: %v", err)
		return
	}

	logger := log.New(f, "", log.LstdFlags)
	g.engine.crawlers[key] = &registeredCrawler{
		path:    groupPath,
		name:    name,
		crawler: crawler,
		log:     logger,
	}
	g.engine.logFiles[key] = f
}

func (g *CrawlerGroup) Group(name string) *CrawlerGroup {
	name = strings.TrimSpace(name)
	if name == "" {
		return g
	}

	g.engine.mu.Lock()
	defer g.engine.mu.Unlock()

	if child, ok := g.children[name]; ok {
		return child
	}

	newPath := filepath.ToSlash(filepath.Join(cleanGroupPath(g.path), name))
	newGroup := &CrawlerGroup{
		path:     newPath,
		engine:   g.engine,
		parent:   g,
		children: make(map[string]*CrawlerGroup),
	}
	g.children[name] = newGroup

	groupDir := filepath.Join(g.engine.logRootDir, filepath.FromSlash(newPath))
	if err := os.MkdirAll(groupDir, 0o755); err != nil {
		g.engine.logf(newPath, "group", "创建目录失败: %v", err)
	}
	return newGroup
}

func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var firstErr error
	for key, f := range e.logFiles {
		if f == nil {
			continue
		}
		if err := f.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(e.logFiles, key)
	}
	return firstErr
}

func (e *Engine) crawlerLogPath(groupPath, crawlerName string) string {
	groupPath = cleanGroupPath(groupPath)
	fileName := sanitizeFileName(crawlerName) + ".log"
	return filepath.Join(e.logRootDir, filepath.FromSlash(groupPath), fileName)
}

func cleanGroupPath(path string) string {
	p := strings.TrimSpace(path)
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.Trim(p, "/")
	p = strings.ReplaceAll(p, "..", "_")
	if p == "" {
		return "root"
	}
	return p
}

func sanitizeFileName(name string) string {
	n := strings.TrimSpace(name)
	if n == "" {
		return "unknown"
	}
	r := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return r.Replace(n)
}

func (e *Engine) logf(path, name, format string, args ...any) {
	e.globalLogger.Printf("[%s][%s] %s", cleanGroupPath(path), name, fmt.Sprintf(format, args...))
}
