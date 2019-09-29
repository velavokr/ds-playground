package main

import (
	demo "github.com/velavokr/dsplayground/demoserver"
	"github.com/velavokr/dsplayground/demoserver/nodeenv"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/network"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/timer"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/fifoperfect"
)

func main() {
	env := runner.InitFromCommandLine()
	nodeEnv := nodeenv.NewNodeEnv(env,
		network.NewFairLossTcp,
		timer.NewTimer,
	)
	netHandler := link.NewDemoLinkReceiver(env)
	net := fifoperfect.NewFifoPerfectLink(netHandler, nodeEnv)
	reqHandler := link.NewDemoLinkSender(env, net)
	demo.RunServer(env, reqHandler)
}
