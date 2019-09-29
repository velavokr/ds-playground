package main

import (
	"github.com/velavokr/dsplayground/demoserver/nodeenv/network"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/timer"
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/stubborn"
)

func main() {
	link.RunLinkDemo(
		stubborn.NewStubbornLink,
		network.NewFairLossTcp,
		timer.NewTimer,
	)
}
