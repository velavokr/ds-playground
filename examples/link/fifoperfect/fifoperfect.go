package fifoperfect

import (
	"encoding/binary"
	"fmt"
	"github.com/velavokr/gdaf"
)

func NewFifoPerfectLink(handler gdaf.NetHandler, env gdaf.NodeEnv) gdaf.Net {
	pl := &perfectFifoLink{
		handler:   handler,
		toSend:    map[gdaf.NodeName]*sendCtx{},
		toDeliver: map[gdaf.NodeName]*uint64{},
	}
	pl.timer = env.Timer(pl)
	pl.fairLossLink = env.Net(pl)
	return pl
}

type perfectFifoLink struct {
	timer        gdaf.Timer
	handler      gdaf.NetHandler
	fairLossLink gdaf.Net
	toSend       map[gdaf.NodeName]*sendCtx
	toDeliver    map[gdaf.NodeName]*uint64
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

func (pfl *perfectFifoLink) SendMessage(dst gdaf.NodeName, rawMsg []byte) {
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

func (pfl *perfectFifoLink) ReceiveMessage(src gdaf.NodeName, rawMsg []byte) {
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

func (pfl perfectFifoLink) HandleTimer(dst interface{}, id gdaf.TimerId) {
	ctx := pfl.toSend[dst.(gdaf.NodeName)]
	for _, msg := range ctx.msgs {
		pfl.fairLossLink.SendMessage(dst.(gdaf.NodeName), msg)
	}
}

func (pfl perfectFifoLink) retransmit(dst gdaf.NodeName, ctx sendCtx) {
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
