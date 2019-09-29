package besteffort

import (
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/examples/bcast"
	"github.com/velavokr/gdaf/examples/link/perfect"
)

func NewBestEffortBroadcastNet(group gdaf.Group, handler gdaf.NetHandler, env gdaf.NodeEnv) bcast.BroadcastNet {
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

func (b *bestEffort) SendMessage(dst gdaf.NodeName, message []byte) {
	b.perfect.SendMessage(dst, message)
}

func (b *bestEffort) ReceiveMessage(src gdaf.NodeName, message []byte) {
	b.handler.ReceiveMessage(src, message)
}

type bestEffort struct {
	handler gdaf.NetHandler
	perfect gdaf.Net
	group   gdaf.Group
}