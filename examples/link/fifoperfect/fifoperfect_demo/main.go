package main

import (
	demo "github.com/velavokr/gdaf/demoserver"
	"github.com/velavokr/gdaf/demoserver/nodeenv"
	"github.com/velavokr/gdaf/demoserver/nodeenv/network"
	"github.com/velavokr/gdaf/demoserver/nodeenv/timer"
	"github.com/velavokr/gdaf/demoserver/runner"
	"github.com/velavokr/gdaf/examples/link"
	"github.com/velavokr/gdaf/examples/link/fifoperfect"
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
