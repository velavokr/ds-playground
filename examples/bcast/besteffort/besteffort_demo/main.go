package main

import (
	"github.com/velavokr/dsplayground/examples/bcast"
	"github.com/velavokr/dsplayground/examples/bcast/besteffort"
)

func main() {
	bcast.RunBcastDemo(besteffort.NewBestEffortBroadcastNet)
}
