package link

import (
	"github.com/velavokr/dsplayground/ifaces"
)

type Link interface {
	ifaces.Net
	ifaces.NetHandler
}

type NewLink = func(handler ifaces.NetHandler, env ifaces.NodeEnv) Link
