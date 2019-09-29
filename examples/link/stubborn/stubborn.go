package stubborn

import (
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/examples/link"
)

func NewStubbornLink(handler gdaf.NetHandler, env gdaf.NodeEnv) link.Link {
	lnk := &stubbornLink{
		handler: handler,
	}
	lnk.timer = env.Timer(lnk)
	lnk.fairLossLink = env.Net(lnk)
	return lnk
}

type stubbornLink struct {
	timer        gdaf.Timer
	handler      gdaf.NetHandler
	fairLossLink gdaf.Net
	toSend       []ctx
}

type ctx struct {
	msg []byte
	dst gdaf.NodeName
}

func (sl *stubbornLink) SendMessage(dst gdaf.NodeName, msg []byte) {
	ctx := ctx{
		msg: msg,
		dst: dst,
	}
	sl.toSend = append(sl.toSend, ctx)
	sl.fairLossLink.SendMessage(dst, msg)
	sl.timer.NextTick(ctx)
}

func (sl *stubbornLink) ReceiveMessage(src gdaf.NodeName, rawMsg []byte) {
	sl.handler.ReceiveMessage(src, rawMsg)
}

func (sl *stubbornLink) HandleTimer(c interface{}, id gdaf.TimerId) {
	ctx := c.(ctx)
	sl.fairLossLink.SendMessage(ctx.dst, ctx.msg)
	sl.timer.NextTick(ctx)
}
