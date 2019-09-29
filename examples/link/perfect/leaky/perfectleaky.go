package leaky

import (
	"encoding/binary"
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/examples/link"
	"github.com/velavokr/gdaf/examples/link/stubborn"
)

func NewPerfectLinkLeaky(handler gdaf.NetHandler, env gdaf.NodeEnv) link.Link {
	lnk := &perfectLinkLeaky{
		handler:   handler,
		delivered: map[frameId]bool{},
	}
	lnk.stubbornLink = stubborn.NewStubbornLink(lnk, env)
	return lnk
}

type perfectLinkLeaky struct {
	handler      gdaf.NetHandler
	stubbornLink gdaf.Net
	delivered    map[frameId]bool
	cnt          uint64
}

type frameId struct {
	src gdaf.NodeName
	seq uint64
}

type frame struct {
	seq uint64
	msg []byte
}

func (pl *perfectLinkLeaky) SendMessage(dst gdaf.NodeName, msg []byte) {
	pl.cnt += 1
	pl.stubbornLink.SendMessage(
		dst,
		encodeMsg(frame{
			seq: pl.cnt,
			msg: msg,
		}),
	)
}

func (pl *perfectLinkLeaky) ReceiveMessage(src gdaf.NodeName, rawMsg []byte) {
	msg := decodeMsg(rawMsg)
	id := frameId{seq: msg.seq, src: src}
	_, ok := pl.delivered[id]
	if !ok {
		pl.delivered[id] = true
		pl.handler.ReceiveMessage(src, msg.msg)
	}
}

func encodeMsg(msg frame) []byte {
	buf := make([]byte, 8+len(msg.msg))
	binary.LittleEndian.PutUint64(buf[0:8], msg.seq)
	if msg.msg != nil {
		copy(buf[8:len(msg.msg)], msg.msg)
	}
	return buf
}

func decodeMsg(msg []byte) frame {
	return frame{
		seq: binary.LittleEndian.Uint64(msg[0:8]),
		msg: msg[8:],
	}
}
