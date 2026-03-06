package easecrawler

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	CrawlerGroup
	log              *log.Logger
	mu               sync.RWMutex
	crawlers         map[string]*registeredCrawler
	groupLoggers     map[string]*log.Logger
	groupLogFiles    map[string]*os.File
	autoGroupLogger  bool
	autoLogDirectory string
}

type registeredCrawler struct {
	path    string
	name    string
	crawler Crawler
}

type CrawlerGroup struct {
	Path   string
	Engine *Engine
}

func NewEngine() *Engine {
	e := &Engine{
		CrawlerGroup:     CrawlerGroup{Path: "/"},
		log:              log.New(os.Stdout, "", log.LstdFlags),
		crawlers:         make(map[string]*registeredCrawler),
		groupLoggers:     make(map[string]*log.Logger),
		groupLogFiles:    make(map[string]*os.File),
		autoGroupLogger:  true,
		autoLogDirectory: "logs",
	}
	e.Engine = e
	return e
}

func (e *Engine) SetGroupLogger(path string, logger *log.Logger) {
	path = normalizePath(path)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.groupLoggers[path] = logger
}

func (e *Engine) SetAutoGroupLogger(enable bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.autoGroupLogger = enable
}

func (e *Engine) SetAutoLogDirectory(dir string) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = "logs"
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.autoLogDirectory = dir
}

func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	var firstErr error
	for key, f := range e.groupLogFiles {
		if f == nil {
			continue
		}
		if err := f.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(e.groupLogFiles, key)
	}
	return firstErr
}

func (g *CrawlerGroup) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	g.RunWithContext(context.Background())
}

func (g *CrawlerGroup) RunWithContext(ctx context.Context) {
	g.Engine.mu.RLock()
	items := make([]*registeredCrawler, 0, len(g.Engine.crawlers))
	for _, c := range g.Engine.crawlers {
		items = append(items, c)
	}
	g.Engine.mu.RUnlock()

	for _, item := range items {
		go g.Engine.runCrawlerLoop(ctx, item)
	}
}

func (e *Engine) runCrawlerLoop(ctx context.Context, item *registeredCrawler) {
	meta := item.crawler.Meta()
	interval := meta.Interval
	if interval <= 0 {
		interval = time.Second
	}

	runOnce := func() {
		logger := e.crawlerLogger(item.path, item.name)
		cctx := &Context{}
		cctx.Set(ContextLoggerKey, logger)
		err := item.crawler.Run(ctx, cctx)
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
	path := normalizePath(g.Path)
	if metaPath := normalizePath(crawler.Meta().Group); metaPath != "" {
		path = metaPath
	}
	name := crawler.Name()
	key := joinPath(path, name)

	g.Engine.mu.Lock()
	defer g.Engine.mu.Unlock()
	g.Engine.crawlers[key] = &registeredCrawler{path: path, name: name, crawler: crawler}
}

func (g *CrawlerGroup) Group(name string) *CrawlerGroup {
	return &CrawlerGroup{Path: joinPath(g.Path, name), Engine: g.Engine}
}

func (e *Engine) logf(path, crawler, format string, args ...any) {
	prefix := "[path:" + normalizePath(path) + "][crawler:" + crawler + "] "
	e.pickLogger(path).Printf(prefix+format, args...)
}

func (e *Engine) crawlerLogger(path, crawler string) *log.Logger {
	base := e.pickLogger(path)
	prefix := "[path:" + normalizePath(path) + "][crawler:" + crawler + "] "
	return log.New(base.Writer(), prefix, log.LstdFlags)
}

func (e *Engine) pickLogger(path string) *log.Logger {
	path = normalizePath(path)

	e.mu.RLock()
	if l, ok := e.groupLoggers[path]; ok && l != nil {
		e.mu.RUnlock()
		return l
	}

	bestLen := -1
	var best *log.Logger
	for p, l := range e.groupLoggers {
		if l == nil {
			continue
		}
		p = normalizePath(p)
		if path == p || strings.HasPrefix(path, p+"/") {
			if len(p) > bestLen {
				bestLen = len(p)
				best = l
			}
		}
	}
	autoEnabled := e.autoGroupLogger
	logDir := e.autoLogDirectory
	e.mu.RUnlock()

	if best != nil {
		return best
	}
	if autoEnabled {
		if l := e.ensureAutoGroupLogger(path, logDir); l != nil {
			return l
		}
	}
	return e.log
}

func (e *Engine) ensureAutoGroupLogger(path, dir string) *log.Logger {
	e.mu.Lock()
	defer e.mu.Unlock()

	if l, ok := e.groupLoggers[path]; ok && l != nil {
		return l
	}
	if dir == "" {
		dir = "logs"
	}

	path = normalizePath(path)
	groupDir := filepath.Join(dir, filepath.FromSlash(path))
	if err := os.MkdirAll(groupDir, 0o755); err != nil {
		e.log.Printf("[path:%s][crawler:system] create log dir failed: %v", path, err)
		return nil
	}

	fullPath := filepath.Join(groupDir, "current.log")
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		e.log.Printf("[path:%s][crawler:system] open log file failed: %v", path, err)
		return nil
	}

	l := log.New(f, "", log.LstdFlags)
	e.groupLoggers[path] = l
	e.groupLogFiles[path] = f
	return l
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.Trim(path, "/")
	if path == "" {
		return "root"
	}
	path = strings.ReplaceAll(path, "..", "_")
	return path
}

func joinPath(parent, child string) string {
	parent = strings.Trim(strings.ReplaceAll(parent, "\\", "/"), "/")
	child = strings.Trim(strings.ReplaceAll(child, "\\", "/"), "/")

	switch {
	case parent == "" && child == "":
		return "root"
	case parent == "":
		return child
	case child == "":
		return parent
	default:
		return parent + "/" + child
	}
}
