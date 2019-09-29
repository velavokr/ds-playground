package fifoperfect

import (
	"encoding/binary"
	"fmt"
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/ifaces"
)

func NewFifoPerfectLink(handler ifaces.NetHandler, env ifaces.NodeEnv) link.Link {
	pl := &perfectFifoLink{
		handler:   handler,
		toSend:    map[ifaces.NodeName]*sendCtx{},
		toDeliver: map[ifaces.NodeName]*uint64{},
	}
	pl.timer = env.Timer(pl)
	pl.fairLossLink = env.Net(pl)
	return pl
}

type perfectFifoLink struct {
	timer        ifaces.Timer
	handler      ifaces.NetHandler
	fairLossLink ifaces.Net
	toSend       map[ifaces.NodeName]*sendCtx
	toDeliver    map[ifaces.NodeName]*uint64
}

const (
	typeSnd = 0
	typeAck = 1
)

type sendCtx struct {
	msgs [][]byte
	seq  uint64
}

type frame struct {
	typeId uint32
	seq    uint64
	msg    []byte
}

func (pfl *perfectFifoLink) SendMessage(dst ifaces.NodeName, rawMsg []byte) {
	nxt, ok := pfl.toSend[dst]
	if !ok {
		msgs := make([][]byte, 0, 1)
		var seq uint64
		nxt = &sendCtx{
			msgs: msgs,
			seq:  seq,
		}
		pfl.toSend[dst] = nxt
	}
	nxt.msgs = append(nxt.msgs, encodeMsg(frame{typeId: typeSnd, seq: nxt.seq, msg: rawMsg}))
	nxt.seq += 1
	pfl.fairLossLink.SendMessage(dst, nxt.msgs[0])
	pfl.timer.NextTick(dst)
}

func (pfl *perfectFifoLink) ReceiveMessage(src ifaces.NodeName, rawMsg []byte) {
	msg := decodeMsg(rawMsg)
	nxt, ok := pfl.toDeliver[src]
	if !ok {
		pfl.toDeliver[src] = new(uint64)
		nxt = pfl.toDeliver[src]
	}
	switch msg.typeId {
	case typeSnd:
		pfl.SendMessage(src, encodeMsg(frame{typeId: typeAck, seq: *nxt}))
		if msg.seq == *nxt {
			pfl.handler.ReceiveMessage(src, msg.msg)
			*nxt += 1
		}
		pfl.fairLossLink.SendMessage(
			src,
			encodeMsg(frame{
				typeId: typeAck,
				seq:    *nxt,
			}),
		)
	case typeAck:
		ctx, ok := pfl.toSend[src]
		if !ok {
			panic(fmt.Sprintf("unexpected ack"))
		}
		if msg.seq >= ctx.seq {
			panic(fmt.Sprintf("unexpected ack num %d", msg.seq))
		}
		beg := ctx.seq - uint64(len(ctx.msgs))
		if msg.seq >= beg {
			ctx.msgs = (ctx.msgs)[msg.seq-beg+1:]
		}
		if msg.seq+1 != ctx.seq {
			pfl.retransmit(src, *ctx)
		}
	default:
		panic(fmt.Sprintf("unexpected typeId %d", msg.typeId))
	}
}

func (pfl perfectFifoLink) HandleTimer(dst interface{}, id ifaces.TimerId) {
	ctx := pfl.toSend[dst.(ifaces.NodeName)]
	for _, msg := range ctx.msgs {
		pfl.fairLossLink.SendMessage(dst.(ifaces.NodeName), msg)
	}
}

func (pfl perfectFifoLink) retransmit(dst ifaces.NodeName, ctx sendCtx) {
	for _, m := range ctx.msgs {
		pfl.fairLossLink.SendMessage(dst, m)
	}
}

func encodeMsg(msg frame) []byte {
	buf := make([]byte, 12+len(msg.msg))
	binary.LittleEndian.PutUint32(buf, msg.typeId)
	binary.LittleEndian.PutUint64(buf, msg.seq)
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
