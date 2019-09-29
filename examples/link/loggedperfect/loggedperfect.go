package loggedperfect

import (
	"encoding/binary"
	"fmt"
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/examples/link"
)

func NewLoggedPerfectLink(handler gdaf.NetHandler, env gdaf.NodeEnv) link.Link {
	lpl := &loggedPerfectLink{
		netHandler: handler,
	}
	lpl.timer = env.Timer(lpl)
	lpl.fairLossLink = env.Net(lpl)
	storage := env.Storage()
	lpl.cnt = storage.OpenTable("cnt")
	lpl.toSend = storage.OpenTable("to_send")
	lpl.delivered = storage.OpenTable("delivered")
	return lpl
}

type loggedPerfectLink struct {
	timer        gdaf.Timer
	netHandler   gdaf.NetHandler
	fairLossLink gdaf.Net
	cnt          gdaf.DiskTable
	toSend       gdaf.DiskTable
	delivered    gdaf.DiskTable
}

const (
	typeSnd  = 0
	typeAck  = 1
	typeAck2 = 2
	cntKey   = "cnt"
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

func (lpl *loggedPerfectLink) SendMessage(dst gdaf.NodeName, rawMsg []byte) {
	cnt := decodeCnt(lpl.cnt.LoadValue([]byte(cntKey))) + 1
	ctx := sendCtx{
		msg: lpl.encodeFrame(frame{
			typeId: typeSnd,
			seq:    cnt,
			msg:    rawMsg,
		}),
		dst: dst,
	}
	cntVal := encodeCnt(cnt)
	lpl.cnt.StoreValue([]byte(cntKey), cntVal)
	lpl.toSend.StoreValue(cntVal, encodeCtx(ctx))
	lpl.fairLossLink.SendMessage(ctx.dst, ctx.msg)
	lpl.timer.NextTick(cnt)
}

func (lpl *loggedPerfectLink) ReceiveMessage(src gdaf.NodeName, rawMsg []byte) {
	msg := lpl.decodeFrame(rawMsg)
	id := encodeId(id{
		src: src,
		seq: msg.seq,
	})
	switch msg.typeId {
	case typeSnd:
		lpl.fairLossLink.SendMessage(
			src,
			lpl.encodeFrame(frame{
				typeId: typeAck,
				seq:    msg.seq,
			}),
		)
		ok := lpl.delivered.LoadValue(id)
		if ok == nil {
			lpl.delivered.StoreValue(id, []byte{})
			lpl.netHandler.ReceiveMessage(src, msg.msg)
		}
	case typeAck:
		lpl.toSend.DeleteKey(encodeCnt(msg.seq))
		lpl.fairLossLink.SendMessage(src, lpl.encodeFrame(frame{
			typeId: typeAck2,
			seq:    msg.seq,
		}))
	case typeAck2:
		lpl.delivered.DeleteKey(id)
	default:
		panic(fmt.Sprintf("unexpected typeId %d", msg.typeId))
	}
}

func (lpl *loggedPerfectLink) HandleTimer(cnt interface{}, id gdaf.TimerId) {
	rawCtx := lpl.toSend.LoadValue(encodeCnt(cnt.(uint64)))
	if rawCtx != nil {
		ctx := decodeCtx(rawCtx)
		lpl.fairLossLink.SendMessage(ctx.dst, ctx.msg)
		lpl.timer.NextTick(cnt)
	}
}

func (lpl loggedPerfectLink) encodeFrame(msg frame) []byte {
	buf := make([]byte, 12+len(msg.msg))
	binary.LittleEndian.PutUint32(buf[0:4], msg.typeId)
	binary.LittleEndian.PutUint64(buf[4:12], msg.seq)
	if msg.msg != nil {
		copy(buf[12:], msg.msg)
	}
	return buf
}

func (lpl loggedPerfectLink) decodeFrame(msg []byte) frame {
	return frame{
		typeId: binary.LittleEndian.Uint32(msg[0:4]),
		seq:    binary.LittleEndian.Uint64(msg[4:12]),
		msg:    msg[12:],
	}
}

func encodeCnt(cnt uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, cnt)
	return buf
}

func decodeCnt(rawCnt []byte) uint64 {
	if rawCnt == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(rawCnt)
}

func encodeId(id id) []byte {
	buf := make([]byte, 8+len(id.src))
	binary.LittleEndian.PutUint64(buf[0:8], id.seq)
	copy(buf[8:], id.src)
	return buf
}

func encodeCtx(ctx sendCtx) []byte {
	ln := len(ctx.dst)
	buf := make([]byte, 4+ln+len(ctx.msg))
	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(ctx.dst)))
	copy(buf[4:4+ln], ctx.dst)
	copy(buf[4+ln:], ctx.msg)
	return buf
}

func decodeCtx(rawCtx []byte) sendCtx {
	ln := binary.LittleEndian.Uint32(rawCtx)
	return sendCtx{
		dst: gdaf.NodeName(rawCtx[4 : 4+ln]),
		msg: rawCtx[4+ln:],
	}
}
