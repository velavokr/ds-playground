package main

import (
	"github.com/velavokr/gdaf/demoserver/nodeenv/network"
	"github.com/velavokr/gdaf/demoserver/nodeenv/storage"
	"github.com/velavokr/gdaf/demoserver/nodeenv/timer"
	"github.com/velavokr/gdaf/examples/link"
	"github.com/velavokr/gdaf/examples/link/loggedperfect"
)

func main() {
	link.RunLinkDemo(
		loggedperfect.NewLoggedPerfectLink,
		network.NewFairLossTcp,
		timer.NewTimer,
		storage.NewStorage,
	)
}
