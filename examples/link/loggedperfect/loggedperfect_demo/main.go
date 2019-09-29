package main

import (
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/loggedperfect"
)

func main() {
	link.RunLinkDemo(loggedperfect.NewLoggedPerfectLink)
}
