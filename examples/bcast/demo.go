package bcast

import (
	"github.com/velavokr/dsplayground/ifaces"
	demo "github.com/velavokr/dsplayground/demoserver"
	"github.com/velavokr/dsplayground/demoserver/nodeenv"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/demoserver/utils"
	"net/url"
)

func RunBcastDemo(newBcast NewBroadcastNet, makers ...interface{}) {
	rt := runner.InitFromCommandLine()
	nodeEnv := nodeenv.NewNodeEnv(rt, makers...)
	netHandler := &linkReceiver{rt: rt,}
	link := newBcast(rt.Cfg.Group, netHandler, nodeEnv)
	reqHandler := &linkSender{rt: rt, link: link,}
	demo.RunServer(rt, reqHandler)
}

func (h *linkSender) HandleApiCall(url *url.URL, b []byte) ([]byte, error) {
	q := url.Query()
	dst := q.Get("dst")
	msg := q.Get("msg")
	if dst == "bcast" {
		h.rt.Run(func() {
			h.link.Broadcast([]byte(msg))
		}, runner.ExitOnPanic, "sending to ", dst, msg)
	} else {
		h.rt.Run(func() {
			h.link.SendMessage(dst, []byte(msg))
		}, runner.ExitOnPanic, "sending to ", dst, msg)
	}
	return nil, nil
}

func (h *linkReceiver) ReceiveMessage(src ifaces.NodeName, msg []byte) {
	h.rt.Println(false, utils.Sprint("delivered from", src, msg))
}

type linkReceiver struct {
	rt *runner.Runtime
}

type linkSender struct {
	rt   *runner.Runtime
	link BroadcastNet
}
