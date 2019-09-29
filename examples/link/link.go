package link

import (
	"github.com/velavokr/gdaf"
)

type Link interface {
	gdaf.Net
	gdaf.NetHandler
}

type NewLink = func(handler gdaf.NetHandler, env gdaf.NodeEnv) Link
