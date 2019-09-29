package timer

import (
	"context"
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/demoserver/runner"
	"time"
)

func NewTimer(cfg *runner.Runtime, handler gdaf.TimerHandler) gdaf.Timer {
	return &timer{
		rt:      cfg,
		alarms:  alarmsMap{},
		handler: handler,
	}
}

type alarmsMap map[gdaf.TimerId]context.CancelFunc

type timer struct {
	rt      *runner.Runtime
	handler gdaf.TimerHandler
	alarms  alarmsMap
	cnt     gdaf.TimerId
}

func (t *timer) After(ticks uint32, alarmCtx interface{}) gdaf.TimerId {
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

func (t *timer) NextTick(ctx interface{}) gdaf.TimerId {
	return t.After(1, ctx)
}

func (t *timer) CancelTimer(id gdaf.TimerId) {
	// The call is guarded by the RunGuarded mutex
	t.rt.Run(func() {
		cancel := t.alarms[id]
		if cancel != nil {
			cancel()
		}
	}, runner.ExitOnPanic|runner.VerboseLog, "cancel timer", id)
}
