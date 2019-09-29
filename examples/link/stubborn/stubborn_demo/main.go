package main

import (
	"github.com/velavokr/dsplayground/examples/link"
	"github.com/velavokr/dsplayground/examples/link/stubborn"
)

func main() {
	link.RunLinkDemo(stubborn.NewStubbornLink)
}
