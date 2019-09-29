package main

import (
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/perfect"
)

func main() {
	link.RunLinkDemo(perfect.NewPerfectLink)
}
