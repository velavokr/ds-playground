package nodeenv

import (
	"github.com/velavokr/dsplayground/demoserver/nodeenv/network"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/storage"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/timer"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/ifaces"
)

type NetMaker = func(rt *runner.Runtime, handler ifaces.NetHandler) ifaces.Net
type TimerMaker = func(rt *runner.Runtime, handler ifaces.TimerHandler) ifaces.Timer
type StorageMaker = func(rt *runner.Runtime) ifaces.Storage

func NewNodeEnv(rt *runner.Runtime) ifaces.NodeEnv {
	return &nodeEnv{rt: rt,}
}

type nodeEnv struct {
	rt      *runner.Runtime
	net     NetMaker
	timer   TimerMaker
	storage StorageMaker
}

func (n *nodeEnv) Net(handler ifaces.NetHandler) ifaces.Net {
	n.rt.Println(false, "initialize net")
	return network.NewFairLossTcp(n.rt, handler)
}

func (n *nodeEnv) Timer(handler ifaces.TimerHandler) ifaces.Timer {
	n.rt.Println(false, "initialize timer")
	return timer.NewTimer(n.rt, handler)
}

func (n nodeEnv) Storage() ifaces.Storage {
	n.rt.Println(false, "initialize storage")
	return storage.NewStorage(n.rt)
}

func (n nodeEnv) PKI() ifaces.PKI {
	panic("implement me")
}
