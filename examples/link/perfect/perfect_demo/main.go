package main

import (
	"github.com/velavokr/dsplayground/demoserver/nodeenv/network"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/timer"
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/perfect"
)

func main() {
	link.RunLinkDemo(
		perfect.NewPerfectLink,
		network.NewFairLossTcp,
		timer.NewTimer,
	)
}
