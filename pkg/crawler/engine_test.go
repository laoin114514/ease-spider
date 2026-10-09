package crawler

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// stubCrawler 是调度测试用的最小插件实现。
type stubCrawler struct {
	name             string
	interval         time.Duration
	startImmediately bool
	// run 是插件体；为空时只累加执行轮数
	run func(c *Context) error
	// runs 记录被执行了多少轮
	runs atomic.Int64
}

func (s *stubCrawler) Name() string { return s.name }

func (s *stubCrawler) Meta() Meta {
	return Meta{Interval: s.interval, StartImmediately: s.startImmediately}
}

func (s *stubCrawler) Run(c *Context) error {
	s.runs.Add(1)
	if s.run != nil {
		return s.run(c)
	}
	return nil
}

// newTestEngine 构造一个日志写到临时目录的引擎，避免污染仓库里的 logs/。
func newTestEngine(t *testing.T) *Engine {
	t.Helper()
	e := New()
	e.SetLogRootDir(t.TempDir())
	// 引擎持有插件日志文件句柄，测试结束前必须释放：
	// 否则 Windows 上 t.TempDir() 的清理会因文件被占用而失败
	t.Cleanup(func() { _ = e.Close() })
	return e
}

// waitReturned 等待 RunWithContext 返回，超时即判定失败。
func waitReturned(t *testing.T, returned <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-returned:
	case <-time.After(3 * time.Second):
		t.Fatalf("%s: RunWithContext 未在 3s 内返回", what)
	}
}

// TestRunWithContextStopsOnCancel 验证 ctx 取消后调度循环退出、不再执行新的一轮。
func TestRunWithContextStopsOnCancel(t *testing.T) {
	e := newTestEngine(t)

	started := make(chan struct{})
	var once sync.Once
	c := &stubCrawler{
		name:             "idle",
		interval:         time.Hour,
		startImmediately: true,
		run: func(*Context) error {
			once.Do(func() { close(started) })
			return nil
		},
	}
	e.Group("g").Register(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	returned := make(chan struct{})
	go func() {
		e.RunWithContext(ctx)
		close(returned)
	}()

	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("启动即跑没有执行")
	}

	cancel()
	waitReturned(t, returned, "ctx 取消后")

	if got := c.runs.Load(); got != 1 {
		t.Fatalf("插件执行轮数 = %d, want 1（取消后不应再跑）", got)
	}
}

// TestRunWithContextWaitsForRunningTask 验证收到退出信号时不会打断正在执行的任务。
func TestRunWithContextWaitsForRunningTask(t *testing.T) {
	e := newTestEngine(t)
	e.SetShutdownTimeout(10 * time.Second)

	started := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })

	c := &stubCrawler{
		name:             "slow",
		interval:         time.Hour,
		startImmediately: true,
		run: func(*Context) error {
			close(started)
			<-release
			return nil
		},
	}
	e.Group("g").Register(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	returned := make(chan struct{})
	go func() {
		e.RunWithContext(ctx)
		close(returned)
	}()

	<-started
	cancel()

	// 任务还卡在 run 里，此时不应该返回
	select {
	case <-returned:
		t.Fatal("任务尚未结束，RunWithContext 就返回了（收尾等待没有生效）")
	case <-time.After(200 * time.Millisecond):
	}

	// 放开任务后应当很快收尾返回
	releaseOnce.Do(func() { close(release) })
	waitReturned(t, returned, "任务收尾后")
}

// TestRunWithContextForceExitsAfterShutdownTimeout 验证收尾超时后引擎不再无限期等待。
func TestRunWithContextForceExitsAfterShutdownTimeout(t *testing.T) {
	e := newTestEngine(t)
	e.SetShutdownTimeout(100 * time.Millisecond)

	started := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })

	c := &stubCrawler{
		name:             "stuck",
		interval:         time.Hour,
		startImmediately: true,
		run: func(*Context) error {
			close(started)
			<-release
			return nil
		},
	}
	e.Group("g").Register(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	returned := make(chan struct{})
	go func() {
		e.RunWithContext(ctx)
		close(returned)
	}()

	<-started
	cancel()

	waitReturned(t, returned, "收尾超时后")
}

// TestRunWithContextSkipsStartImmediatelyWhenAlreadyCancelled 验证 ctx 已取消时不再启动新任务。
func TestRunWithContextSkipsStartImmediatelyWhenAlreadyCancelled(t *testing.T) {
	e := newTestEngine(t)
	c := &stubCrawler{name: "idle", interval: time.Hour, startImmediately: true}
	e.Group("g").Register(c)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	returned := make(chan struct{})
	go func() {
		e.RunWithContext(ctx)
		close(returned)
	}()

	waitReturned(t, returned, "ctx 启动前已取消时")

	if got := c.runs.Load(); got != 0 {
		t.Fatalf("ctx 已取消却仍执行了 %d 轮任务, want 0", got)
	}
}

// TestRunWithContextWithoutCrawlersReturns 验证没有注册插件时不会阻塞。
func TestRunWithContextWithoutCrawlersReturns(t *testing.T) {
	e := newTestEngine(t)

	returned := make(chan struct{})
	go func() {
		e.RunWithContext(context.Background())
		close(returned)
	}()

	waitReturned(t, returned, "没有插件时")
}

// TestCloseReleasesLogFilesAndIsIdempotent 验证 Close 释放日志文件且可重复调用。
func TestCloseReleasesLogFilesAndIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	e := New()
	e.SetLogRootDir(dir)
	e.Group("g").Register(&stubCrawler{name: "x", interval: time.Minute})

	logPath := filepath.Join(dir, "g", "x.log")
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("注册后未创建插件日志文件: %v", err)
	}

	if err := e.Close(); err != nil {
		t.Fatalf("首次 Close() 出错: %v", err)
	}
	if err := e.Close(); err != nil {
		t.Fatalf("重复 Close() 出错: %v", err)
	}
}
