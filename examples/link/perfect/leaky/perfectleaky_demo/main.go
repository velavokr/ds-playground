package main

import (
	"github.com/velavokr/gdaf/demoserver/nodeenv/network"
	"github.com/velavokr/gdaf/demoserver/nodeenv/timer"
	"github.com/velavokr/gdaf/examples/link"
	"github.com/velavokr/gdaf/examples/link/perfect/leaky"
)

func main() {
	link.RunLinkDemo(
		leaky.NewPerfectLinkLeaky,
		network.NewFairLossTcp,
		timer.NewTimer,
	)
}
