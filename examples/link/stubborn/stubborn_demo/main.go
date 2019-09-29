package main

import (
	"github.com/velavokr/gdaf/demoserver/nodeenv/network"
	"github.com/velavokr/gdaf/demoserver/nodeenv/timer"
	"github.com/velavokr/gdaf/examples/link"
	"github.com/velavokr/gdaf/examples/link/stubborn"
)

func main() {
	link.RunLinkDemo(
		stubborn.NewStubbornLink,
		network.NewFairLossTcp,
		timer.NewTimer,
	)
}
