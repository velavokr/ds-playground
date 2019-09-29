package network

import (
	"bytes"
	"fmt"
	"github.com/velavokr/dsplayground/demoserver/nodeenv"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/demoserver/utils"
	"github.com/velavokr/dsplayground/ifaces"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestFairLossTcp(t *testing.T) {
	nodes := genNodes()
	envs := genEnvs(nodes)
	handlers := [2]netHandler{}
	nets := [2]ifaces.Net{
		nodeenv.NewNodeEnv(envs[0], NewFairLossTcp).Net(&handlers[0]),
		nodeenv.NewNodeEnv(envs[1], NewFairLossTcp).Net(&handlers[1]),
	}
	time.Sleep(300 * time.Millisecond)
	for from := 0; from < 2; from++ {
		for to := 0; to < 2; to++ {
			nets[from].SendMessage(nodes[to], []byte(fmt.Sprintf("msg%d-%d", from, to)))
		}
	}
	time.Sleep(300 * time.Millisecond)
	envs[0].Cancel()
	envs[0].WaitAll()
	envs[1].Cancel()
	envs[1].WaitAll()
	for to := 0; to < 2; to++ {
		if len(handlers[to].res) != 2 {
			t.Error(handlers[to].res)
		}
		for from := 0; from < 2; from++ {
			checkMsg(t, handlers[to].res[from], nodes[from], fmt.Sprintf("msg%d-%d", from, to))
		}
	}
}

func TestFairLossTcpDrop(t *testing.T) {
	nodes := genNodes()
	envs := genEnvs(nodes)
	handlers := [2]netHandler{}
	nets := [2]ifaces.Net{
		NewFairLossTcp(envs[0], &handlers[0]),
		NewFairLossTcp(envs[1], &handlers[1]),
	}
	time.Sleep(300 * time.Millisecond)
	nets[1].SendMessage(nodes[0], []byte("msg1-0"))
	nets[1].SendMessage(nodes[1], []byte("msg1-1"))
	time.Sleep(30 * time.Millisecond)
	envs[1].Cancel()
	envs[1].WaitAll()
	nets[0].SendMessage(nodes[1], []byte("msg0-1"))
	nets[0].SendMessage(nodes[0], []byte("msg0-0"))
	time.Sleep(300 * time.Millisecond)
	envs[0].Cancel()
	envs[0].WaitAll()
	if len(handlers[0].res) != 2 {
		t.Error(handlers[0].res)
	}
	for i := 0; i < 2; i++ {
		checkMsg(t, handlers[0].res[i], nodes[i], fmt.Sprintf("msg%d-0", i))
	}
	if len(handlers[1].res) != 1 {
		t.Error(handlers[1].res)
	}
	checkMsg(t, handlers[1].res[0], nodes[1], "msg1-1")
}

type netHandler struct {
	res []message
}

func (n *netHandler) ReceiveMessage(src ifaces.NodeName, msg []byte) {
	n.res = append(n.res, message{src: src, data: msg})
	sort.Slice(n.res, func(i, j int) bool {
		return bytes.Compare(n.res[i].data, n.res[j].data) < 0
	})
}

func genNodes() [2]string {
	nodes := [2]string{}
	for i := 0; i < 2; i++ {
		nodes[i] = fmt.Sprintf("127.0.0.1:%d", utils.RandomFreePort())
	}
	return nodes
}

func genEnvs(nodes [2]string) [2]*runner.Runtime {
	b := [2]bytes.Buffer{}
	cfgs := [2]*runner.Runtime{}
	for i := 0; i < 2; i++ {
		cfgs[i] = runner.NewRuntime(runner.UserCfg{
			Group:     ifaces.Group{Nodes: nodes[:], Self: i},
			Tick:      time.Millisecond * 100,
			IoTimeout: time.Second * 10,
		}, &b[i])
	}
	return cfgs
}

func checkMsg(t *testing.T, msg message, src string, data string) {
	if msg.src != src {
		t.Error(msg.src, src)
	}
	if string(msg.data) != data {
		t.Error(msg.data, data)
	}
}

func TestEncodeMsg(t *testing.T) {
	msg := message{
		src:  "ab",
		dst:  "c",
		data: []byte("defgh"),
	}
	b := []byte{2, 0, 0, 0, 'a', 'b', 1, 0, 0, 0, 'c', 5, 0, 0, 0, 'd', 'e', 'f', 'g', 'h'}
	rawMsg := encodeMsg(msg)
	if bytes.Compare(b, rawMsg) != 0 {
		t.Error(b, rawMsg)
	}
	msg1 := decodeMsg(rawMsg)
	if msg1.src != msg.src || msg1.dst != msg.dst || bytes.Compare(msg1.data, msg.data) != 0 {
		t.Errorf("%#v", msg1)
	}
}

func TestDecodeMsg(t *testing.T) {
	expectPanic := func(s string) {
		err := recover()
		if err == nil {
			t.Error()
		}
		if !strings.Contains(err.(error).Error(), s) {
			t.Error()
		}
	}

	t.Run("zero", func(t *testing.T) {
		msg := decodeMsg(make([]byte, 12))
		if msg.src != "" || msg.dst != "" || msg.data == nil || string(msg.data) != "" {
			t.Error()
		}
	})

	t.Run("truncated_zero", func(t *testing.T) {
		defer expectPanic("truncated")
		decodeMsg(make([]byte, 11))
	})

	t.Run("truncated_one", func(t *testing.T) {
		defer expectPanic("truncated")
		b := make([]byte, 12)
		b[0] = 1
		decodeMsg(b)
	})

	t.Run("excess", func(t *testing.T) {
		defer expectPanic("excess")
		decodeMsg(make([]byte, 13))
	})
}
