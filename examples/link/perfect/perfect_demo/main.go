package main

import (
	"github.com/velavokr/gdaf/demoserver/nodeenv/network"
	"github.com/velavokr/gdaf/demoserver/nodeenv/timer"
	"github.com/velavokr/gdaf/examples/link"
	"github.com/velavokr/gdaf/examples/link/perfect"
)

func main() {
	link.RunLinkDemo(
		perfect.NewPerfectLink,
		network.NewFairLossTcp,
		timer.NewTimer,
	)
}
