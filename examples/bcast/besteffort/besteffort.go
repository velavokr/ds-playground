package besteffort

import (
	"github.com/velavokr/dsplayground/ifaces"
	"github.com/velavokr/dsplayground/examples/bcast"
	"github.com/velavokr/dsplayground/examples/link/perfect"
)

func NewBestEffortBroadcastNet(group ifaces.Group, handler ifaces.NetHandler, env ifaces.NodeEnv) bcast.BroadcastNet {
	be := &bestEffort{
		handler: handler,
		group:   group,
	}
	be.perfect = perfect.NewPerfectLink(be, env)
	return be
}

func (b *bestEffort) Broadcast(msg []byte) {
	for _, n := range b.group.Nodes {
		b.SendMessage(n, msg)
	}
}

func (b *bestEffort) SendMessage(dst ifaces.NodeName, message []byte) {
	b.perfect.SendMessage(dst, message)
}

func (b *bestEffort) ReceiveMessage(src ifaces.NodeName, message []byte) {
	b.handler.ReceiveMessage(src, message)
}

type bestEffort struct {
	handler ifaces.NetHandler
	perfect ifaces.Net
	group   ifaces.Group
}