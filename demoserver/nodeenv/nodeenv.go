package nodeenv

import (
	"errors"
	"fmt"
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/demoserver/runner"
)

type NetMaker = func(rt *runner.Runtime, handler gdaf.NetHandler) gdaf.Net
type TimerMaker = func(rt *runner.Runtime, handler gdaf.TimerHandler) gdaf.Timer
type StorageMaker = func(rt *runner.Runtime) gdaf.Storage

func NewNodeEnv(rt *runner.Runtime, makers ...interface{}) gdaf.NodeEnv {
	env := &nodeEnv{rt: rt,}
	for _, m := range makers {
		switch v := m.(type) {
		case NetMaker:
			env.net = v
		case TimerMaker:
			env.timer = v
		case StorageMaker:
			env.storage = v
		default:
			panic(errors.New(fmt.Sprintf("unknown maker %T", v)))
		}
	}
	return env
}

type nodeEnv struct {
	rt      *runner.Runtime
	net     NetMaker
	timer   TimerMaker
	storage StorageMaker
}

func (n *nodeEnv) Net(handler gdaf.NetHandler) gdaf.Net {
	n.rt.Println(false, "initialize net")
	return n.net(n.rt, handler)
}

func (n *nodeEnv) Timer(handler gdaf.TimerHandler) gdaf.Timer {
	n.rt.Println(false, "initialize timer")
	return n.timer(n.rt, handler)
}

func (n nodeEnv) Storage() gdaf.Storage {
	n.rt.Println(false, "initialize storage")
	return n.storage(n.rt)
}

func (n nodeEnv) PKI() gdaf.PKI {
	panic("implement me")
}
