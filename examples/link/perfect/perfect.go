package perfect

import (
	"encoding/binary"
	"fmt"
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/examples/link"
)

func NewPerfectLink(handler gdaf.NetHandler, env gdaf.NodeEnv) link.Link {
	pl := &perfectLink{
		netHandler: handler,
		toSend:     map[uint64]sendCtx{},
		delivered:  map[id]bool{},
	}
	pl.timer = env.Timer(pl)
	pl.fairLossLink = env.Net(pl)
	return pl
}

type perfectLink struct {
	timer        gdaf.Timer
	netHandler   gdaf.NetHandler
	fairLossLink gdaf.Net
	cnt          uint64
	toSend       map[uint64]sendCtx
	delivered    map[id]bool
}

const (
	typeSnd  = 0
	typeAck  = 1
	typeAck2 = 2
)

type id struct {
	src gdaf.NodeName
	seq uint64
}

type sendCtx struct {
	msg []byte
	dst gdaf.NodeName
}

type frame struct {
	typeId uint32
	seq    uint64
	msg    []byte
}

func (pl *perfectLink) SendMessage(dst gdaf.NodeName, rawMsg []byte) {
	pl.cnt += 1
	ctx := sendCtx{
		msg: encodeMsg(frame{
			typeId: typeSnd,
			seq:    pl.cnt,
			msg:    rawMsg,
		}),
		dst: dst,
	}
	pl.toSend[pl.cnt] = ctx
	pl.fairLossLink.SendMessage(ctx.dst, ctx.msg)
	pl.timer.NextTick(pl.cnt)
}

func (pl *perfectLink) ReceiveMessage(src gdaf.NodeName, rawMsg []byte) {
	msg := decodeMsg(rawMsg)
	id := id{
		src: src,
		seq: msg.seq,
	}
	switch msg.typeId {
	case typeSnd:
		pl.fairLossLink.SendMessage(
			src,
			encodeMsg(frame{
				typeId: typeAck,
				seq:    msg.seq,
			}),
		)
		_, ok := pl.delivered[id]
		if !ok {
			pl.delivered[id] = true
			pl.netHandler.ReceiveMessage(src, msg.msg)
		}
	case typeAck:
		delete(pl.toSend, msg.seq)
		pl.fairLossLink.SendMessage(src, encodeMsg(frame{
			typeId: typeAck2,
			seq:    msg.seq,
		}))
	case typeAck2:
		delete(pl.delivered, id)
	default:
		panic(fmt.Sprintf("unexpected typeId %d", msg.typeId))
	}
}

func (pl *perfectLink) HandleTimer(cnt interface{}, id gdaf.TimerId) {
	ctx, ok := pl.toSend[cnt.(uint64)]
	if ok {
		pl.fairLossLink.SendMessage(ctx.dst, ctx.msg)
		pl.timer.NextTick(cnt)
	}
}

func encodeMsg(msg frame) []byte {
	buf := make([]byte, 12+len(msg.msg))
	binary.LittleEndian.PutUint32(buf[0:4], msg.typeId)
	binary.LittleEndian.PutUint64(buf[4:12], msg.seq)
	if msg.msg != nil {
		copy(buf[12:len(msg.msg)], msg.msg)
	}
	return buf
}

func decodeMsg(msg []byte) frame {
	return frame{
		typeId: binary.LittleEndian.Uint32(msg[0:4]),
		seq:    binary.LittleEndian.Uint64(msg[4:12]),
		msg:    msg[12:],
	}
}
