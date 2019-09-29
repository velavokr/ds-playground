package main

import (
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/fifoperfect"
)

func main() {
	link.RunLinkDemo(fifoperfect.NewFifoPerfectLink)
}
