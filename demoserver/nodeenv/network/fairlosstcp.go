package network

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/velavokr/gdaf"
	"github.com/velavokr/gdaf/demoserver/runner"
	"github.com/velavokr/gdaf/demoserver/utils"
	"io"
	"io/ioutil"
	"net"
)

func NewFairLossTcp(env *runner.Runtime, handler gdaf.NetHandler) gdaf.Net {
	addr := env.Cfg.Nodes[env.Cfg.Self]
	env.RunAsync(func(ctx context.Context) {
		listener, err := new(net.ListenConfig).Listen(ctx, "tcp", addr)
		if err != nil {
			panic(err)
		}

		tcpListener := listener.(*net.TCPListener)
		defer closeSocket(tcpListener, env)
		env.Println(true, "tcp listen on ", addr)

		for {
			if err := tcpListener.SetDeadline(utils.ToDeadline(env.Cfg.Tick)); err != nil {
				panic(err)
			}

			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				if ctx.Err() == context.Canceled {
					return
				}
				netErr, ok := err.(net.Error)
				if ok && (netErr.Temporary() || netErr.Timeout()) {
					continue
				}
				panic(err)
			}
			env.Println(true, "tcp accept from ", conn.RemoteAddr())

			env.RunAsync(func(ctx context.Context) {
				defer closeSocket(conn, env)
				if err := setUpConn(conn, env); err != nil {
					panic(err)
				}

				b, err := ioutil.ReadAll(conn)
				if err != nil {
					panic(err)
				}

				msg := decodeMsg(b)
				env.RunGuarded(func() {
					handler.ReceiveMessage(msg.src, msg.data)
				}, runner.ExitOnPanic|runner.VerboseLog, "message delivery ", msg)
			}, runner.VerboseLog, "reading from ", conn.RemoteAddr())
		}
	}, runner.ExitOnPanic|runner.VerboseLog, "listen ", addr)

	return &fairLossNet{
		rt: env,
	}
}

type fairLossNet struct {
	rt *runner.Runtime
}

func (f *fairLossNet) SendMessage(dst gdaf.NodeName, msg []byte) {
	f.rt.RunAsync(func(ctx context.Context) {
		conn, err := (&net.Dialer{Timeout: f.rt.Cfg.IoTimeout}).DialContext(ctx, "tcp", dst)
		if err != nil {
			panic(err)
		}
		defer closeSocket(conn, f.rt)
		tcpConn := conn.(*net.TCPConn)

		if err := setUpConn(tcpConn, f.rt); err != nil {
			panic(err)
		}

		f.rt.Println(true, "connected to", dst)

		msg := message{
			src:  f.rt.Cfg.Nodes[f.rt.Cfg.Self],
			dst:  dst,
			data: msg,
		}

		if err := utils.WriteAll(conn, encodeMsg(msg)); err != nil {
			panic(err)
		}
	}, runner.VerboseLog, "message send ", dst, msg)
}

type message struct {
	src  gdaf.NodeName
	dst  gdaf.NodeName
	data []byte
}

func decodeMsg(rawMsg []byte) message {
	if len(rawMsg) < 3*4 {
		panic(errors.New(fmt.Sprintf("truncated msg (%d)", len(rawMsg))))
	}
	res := [3][]byte{}
	for i := range res {
		if len(rawMsg) < 4 {
			panic(errors.New(fmt.Sprintf("truncated len (%d)", len(rawMsg))))
		}
		ln := int(binary.LittleEndian.Uint32(rawMsg[0:4]))
		rawMsg = rawMsg[4:]
		if len(rawMsg) < ln {
			panic(errors.New(fmt.Sprintf("truncated data (%d)", len(rawMsg))))
		}
		res[i] = rawMsg[:ln]
		rawMsg = rawMsg[ln:]
	}
	if len(rawMsg) > 0 {
		panic(errors.New(fmt.Sprintf("excess tail data (%d)", len(rawMsg))))
	}
	return message{
		src:  gdaf.NodeName(res[0]),
		dst:  gdaf.NodeName(res[1]),
		data: res[2],
	}
}

func encodeMsg(msg message) []byte {
	res := make([]byte, 4+len(msg.src)+4+len(msg.dst)+4+len(msg.data))
	pos := 0
	for _, b := range [][]byte{[]byte(msg.src), []byte(msg.dst), msg.data} {
		binary.LittleEndian.PutUint32(res[pos:], uint32(len(b)))
		pos += 4
		copy(res[pos:], b)
		pos += len(b)
	}
	return res
}

func closeSocket(c io.Closer, rt *runner.Runtime) {
	err := c.Close()
	if err != nil {
		rt.Println(false, err.Error())
	}
}

func setUpConn(conn *net.TCPConn, env *runner.Runtime) error {
	if err := conn.SetNoDelay(true); err != nil {
		return err
	}
	return conn.SetDeadline(utils.ToDeadline(env.Cfg.IoTimeout))
}
