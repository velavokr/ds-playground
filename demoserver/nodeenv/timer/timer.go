package timer

import (
	"context"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/ifaces"
	"time"
)

func NewTimer(cfg *runner.Runtime, handler ifaces.TimerHandler) ifaces.Timer {
	return &timer{
		rt:      cfg,
		alarms:  alarmsMap{},
		handler: handler,
	}
}

type alarmsMap map[ifaces.TimerId]context.CancelFunc

type timer struct {
	rt      *runner.Runtime
	handler ifaces.TimerHandler
	alarms  alarmsMap
	cnt     ifaces.TimerId
}

func (t *timer) After(ticks uint32, alarmCtx interface{}) ifaces.TimerId {
	// The call is guarded by the RunGuarded mutex
	t.cnt += 1
	id := t.cnt
	next := t.rt.Cfg.Tick * time.Duration(ticks)
	t.alarms[id] = t.rt.RunAsyncCancel(func(ctx context.Context) {
		tm := time.NewTimer(next)
		select {
		case _, ok := <-tm.C:
			if ok {
				t.rt.RunGuarded(func() {
					t.handler.HandleTimer(alarmCtx, id)
					delete(t.alarms, id)
				}, runner.ExitOnPanic|runner.VerboseLog, "timer handler ", id, alarmCtx)
			}
		case <-ctx.Done():
			if !tm.Stop() {
				<-tm.C
			}
			t.rt.RunGuarded(func() {
				delete(t.alarms, id)
			}, runner.ExitOnPanic|runner.VerboseLog, "delete timer", id, alarmCtx)
		}
	}, runner.ExitOnPanic|runner.VerboseLog, "timer", id, alarmCtx)
	return id
}

func (t *timer) NextTick(ctx interface{}) ifaces.TimerId {
	return t.After(1, ctx)
}

func (t *timer) CancelTimer(id ifaces.TimerId) {
	// The call is guarded by the RunGuarded mutex
	t.rt.Run(func() {
		cancel := t.alarms[id]
		if cancel != nil {
			cancel()
		}
	}, runner.ExitOnPanic|runner.VerboseLog, "cancel timer", id)
}
