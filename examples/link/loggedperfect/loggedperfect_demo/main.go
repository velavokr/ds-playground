package main

import (
	"github.com/velavokr/dsplayground/demoserver/nodeenv/network"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/storage"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/timer"
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/loggedperfect"
)

func main() {
	link.RunLinkDemo(
		loggedperfect.NewLoggedPerfectLink,
		network.NewFairLossTcp,
		timer.NewTimer,
		storage.NewStorage,
	)
}
