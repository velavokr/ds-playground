package nodeenv

import (
	"errors"
	"fmt"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/ifaces"
)

type NetMaker = func(rt *runner.Runtime, handler ifaces.NetHandler) ifaces.Net
type TimerMaker = func(rt *runner.Runtime, handler ifaces.TimerHandler) ifaces.Timer
type StorageMaker = func(rt *runner.Runtime) ifaces.Storage

func NewNodeEnv(rt *runner.Runtime, makers ...interface{}) ifaces.NodeEnv {
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

func (n *nodeEnv) Net(handler ifaces.NetHandler) ifaces.Net {
	n.rt.Println(false, "initialize net")
	return n.net(n.rt, handler)
}

func (n *nodeEnv) Timer(handler ifaces.TimerHandler) ifaces.Timer {
	n.rt.Println(false, "initialize timer")
	return n.timer(n.rt, handler)
}

func (n nodeEnv) Storage() ifaces.Storage {
	n.rt.Println(false, "initialize storage")
	return n.storage(n.rt)
}

func (n nodeEnv) PKI() ifaces.PKI {
	panic("implement me")
}
