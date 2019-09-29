package main

import (
	"github.com/velavokr/gdaf/demoserver/nodeenv/network"
	"github.com/velavokr/gdaf/demoserver/nodeenv/timer"
	"github.com/velavokr/gdaf/examples/bcast"
	"github.com/velavokr/gdaf/examples/bcast/besteffort"
)

func main() {
	bcast.RunBcastDemo(
		besteffort.NewBestEffortBroadcastNet,
		timer.NewTimer,
		network.NewFairLossTcp,
	)
}
