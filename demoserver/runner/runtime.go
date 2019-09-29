package runner

import (
	"context"
	"fmt"
	"github.com/velavokr/dsplayground/demoserver/logger"
	"github.com/velavokr/dsplayground/demoserver/utils"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
)

type SyncHandler = func()
type AsyncHandler = func(context.Context)

const (
	ExitOnPanic = 1
	VerboseLog  = 2
)

// Runtime manages async jobs and syncronization, signal and panic handling, logging and shutdown
type Runtime struct {
	*logger.Logger
	Cancel context.CancelFunc

	Cfg UserCfg

	ctx        context.Context
	wait       sync.WaitGroup
	handlerMtx sync.Mutex
}

func NewRuntime(cfg UserCfg, logOut io.Writer) *Runtime {
	c, cancel := context.WithCancel(context.Background())
	e := &Runtime{
		Logger: logger.NewLogger(cfg.Verbose, logOut, fmt.Sprintf("[node %d]", cfg.Self)),
		Cancel: cancel,
		Cfg:    cfg,
		ctx:    c,
	}
	e.startSigListener(cancel)
	return e
}

// RunAsync runs the handler in a new goroutine, handles panics and logs the runs.
// Does not prevent data races.
func (c *Runtime) RunAsync(handler AsyncHandler, flags int, comment string, comments ...interface{}) {
	c.doRunAsync(c.Caller(1), c.ctx, handler, flags, comment, comments...)
}

// RunAsyncCancel runs the handler in a new goroutine, handles panics and logs the runs.
// Allows cancelations.
// Does not prevent data races.
func (c *Runtime) RunAsyncCancel(handler AsyncHandler, flags int, comment string, comments ...interface{}) context.CancelFunc {
	ctx, cancel := context.WithCancel(c.ctx)
	c.doRunAsync(c.Caller(1), ctx, handler, flags, comment, comments...)
	return cancel
}

// RunGuarded locks its mutex and runs the handler, handles panics and logs the runs.
// Prevents data races.
// Will deadlock if called recursively.
func (c *Runtime) RunGuarded(handler SyncHandler, flags int, comment string, comments ...interface{}) {
	c.handlerMtx.Lock()
	defer c.handlerMtx.Unlock()
	c.doRun(c.Caller(1), handler, flags, "guarded "+comment, comments...)
}

// Run runs the handler, handles panics and logs the runs.
// Does not prevent data races.
func (c *Runtime) Run(handler SyncHandler, flags int, comment string, comments ...interface{}) {
	c.doRun(c.Caller(1), handler, flags, comment, comments...)
}

// WaitAll wait for all the async jobs to finish
func (c *Runtime) WaitAll() {
	c.wait.Wait()
}

func (c *Runtime) doRunAsync(caller string, ctx context.Context, handler AsyncHandler, flags int, comment string, comments ...interface{}) {
	c.wait.Add(1)
	safeComments := utils.Sprint(comments...)
	go func() {
		defer c.wait.Done()
		c.doRun(caller, func() {
			handler(ctx)
		}, flags, "async "+comment, safeComments)
	}()
}

func (c *Runtime) doRun(caller string, handler SyncHandler, flags int, comment string, comments ...interface{}) {
	defer c.doRecover(caller, (flags&ExitOnPanic == ExitOnPanic) && !c.Cfg.NoCrash, comment)
	safeComments := utils.Sprint(comments...)
	verbose := flags&VerboseLog == VerboseLog
	c.Output(caller, verbose, " run "+comment, safeComments)
	handler()
	c.Output(caller, verbose, "done "+comment, safeComments)
}

func (c *Runtime) doRecover(caller string, exitOnPanic bool, comment string) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			c.Output(caller, false, "panic in", comment, ":", e.Error())
			c.Output(caller, true, string(debug.Stack()))
		default:
			c.Output(caller, false, "panic in", comment, ":", err)
			c.Output(caller, true, string(debug.Stack()))
		}
		if exitOnPanic && !c.Cfg.NoCrash {
			c.Output(caller, false,
				"unexpected panic, further execution does not make sense. Aborting now (use -nocrash to disable this behavior)")
			if !c.Cfg.Verbose {
				c.Output(caller, false, string(debug.Stack()))
			}
			os.Exit(1)
		}
	}
}

func (c *Runtime) startSigListener(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	c.RunAsync(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-sigs:
			cancel()
		}
	}, ExitOnPanic|VerboseLog, "signal listener")
}
