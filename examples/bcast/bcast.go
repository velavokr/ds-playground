package bcast

import (
	"github.com/velavokr/dsplayground/ifaces"
)

type BroadcastNet interface {
	ifaces.Net
	ifaces.NetHandler
	Broadcast(msg []byte)
}

type NewBroadcastNet = func(group ifaces.Group, handler ifaces.NetHandler, env ifaces.NodeEnv) BroadcastNet
