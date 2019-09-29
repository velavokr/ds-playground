package stubborn

import (
	"github.com/velavokr/dsplayground/ifaces"
	"github.com/velavokr/dsplayground/examples/link"
)

func NewStubbornLink(handler ifaces.NetHandler, env ifaces.NodeEnv) link.Link {
	lnk := &stubbornLink{
		handler: handler,
	}
	lnk.timer = env.Timer(lnk)
	lnk.fairLossLink = env.Net(lnk)
	return lnk
}

type stubbornLink struct {
	timer        ifaces.Timer
	handler      ifaces.NetHandler
	fairLossLink ifaces.Net
	toSend       []ctx
}

type ctx struct {
	msg []byte
	dst ifaces.NodeName
}

func (sl *stubbornLink) SendMessage(dst ifaces.NodeName, msg []byte) {
	ctx := ctx{
		msg: msg,
		dst: dst,
	}
	sl.toSend = append(sl.toSend, ctx)
	sl.fairLossLink.SendMessage(dst, msg)
	sl.timer.NextTick(ctx)
}

func (sl *stubbornLink) ReceiveMessage(src ifaces.NodeName, rawMsg []byte) {
	sl.handler.ReceiveMessage(src, rawMsg)
}

func (sl *stubbornLink) HandleTimer(c interface{}, id ifaces.TimerId) {
	ctx := c.(ctx)
	sl.fairLossLink.SendMessage(ctx.dst, ctx.msg)
	sl.timer.NextTick(ctx)
}
