package timer

import (
	"bytes"
	"github.com/velavokr/dsplayground/demoserver/nodeenv"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/ifaces"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	b := bytes.Buffer{}
	rt := runner.NewRuntime(runner.UserCfg{
		Tick:      time.Millisecond * 100,
		IoTimeout: time.Second * 10,
	}, &b)
	h := handler{}
	timer := nodeenv.NewNodeEnv(rt, NewTimer).Timer(&h)

	cnt := 10
	ids := make([]ifaces.TimerId, cnt)
	for i := 0; i < cnt; i++ {
		ids[i] = timer.After(uint32(i), i)
	}
	for i := 0; i < cnt; i += 2 {
		rt.RunGuarded(func() {
			timer.CancelTimer(ids[i])
		}, runner.ExitOnPanic, "cancel timer", ids[i])
	}
	time.Sleep(rt.Cfg.Tick * 12)
	allowed := map[ifaces.TimerId]int{
		1: 0, 2: 1, 4: 3, 6: 5, 8: 7, 10: 9,
	}
	rt.RunGuarded(func() {
		if len(h.alarms) < 5 {
			t.Error(h.alarms)
		} else {
			for _, a := range h.alarms {
				c, ok := allowed[a.id]
				if !ok {
					t.Error(h.alarms)
				}
				if c != a.ctx.(int) {
					t.Error(h.alarms)
				}
			}
		}
	}, runner.ExitOnPanic, "")
	rt.Cancel()
	rt.WaitAll()
}

type alarm struct {
	id  ifaces.TimerId
	ctx interface{}
}

type handler struct {
	alarms []alarm
}

func (h *handler) HandleTimer(ctx interface{}, id ifaces.TimerId) {
	h.alarms = append(h.alarms, alarm{
		id:  id,
		ctx: ctx,
	})
}
