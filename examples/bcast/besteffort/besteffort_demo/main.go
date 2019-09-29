package main

import (
	"github.com/velavokr/dsplayground/demoserver/nodeenv/network"
	"github.com/velavokr/dsplayground/demoserver/nodeenv/timer"
	"github.com/velavokr/dsplayground/examples/bcast"
	"github.com/velavokr/dsplayground/examples/bcast/besteffort"
)

func main() {
	bcast.RunBcastDemo(
		besteffort.NewBestEffortBroadcastNet,
		timer.NewTimer,
		network.NewFairLossTcp,
	)
}
