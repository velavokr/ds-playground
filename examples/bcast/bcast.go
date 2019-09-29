package bcast

import (
	"github.com/velavokr/gdaf"
)

type BroadcastNet interface {
	gdaf.Net
	gdaf.NetHandler
	Broadcast(msg []byte)
}

type NewBroadcastNet = func(group gdaf.Group, handler gdaf.NetHandler, env gdaf.NodeEnv) BroadcastNet
