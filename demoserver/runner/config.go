package runner

import (
	"flag"
	"fmt"
	"github.com/velavokr/dsplayground/ifaces"
	"os"
	"strings"
	"time"
)

type UserCfg struct {
	ifaces.Group
	Http      string
	Tick      time.Duration
	IoTimeout time.Duration
	DbDir     string
	Verbose   bool
	NoCrash   bool
}

func InitFromCommandLine() *Runtime {
	cfg := UserCfg{}

	flag.StringVar(&cfg.Http, "http", "127.0.0.1:8080", "current node http server address")
	var nodesRaw string
	flag.StringVar(&nodesRaw, "nodes", nodesRaw,
		"comma-separated list of all the participating nodes, e.g. '127.0.0.1:8091,127.0.0.1:8092'. All should be different from http addresses")
	flag.IntVar(&cfg.Self, "self", cfg.Self, "current node number")

	flag.DurationVar(&cfg.Tick, "tick", 500*time.Millisecond, "time granularity")
	flag.DurationVar(&cfg.IoTimeout, "iotimeout", 30*time.Second, "network io operations timeout")

	flag.StringVar(&cfg.DbDir, "dbdir", "./db", "directory for storage files")

	flag.BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "verbose logging")
	flag.BoolVar(&cfg.NoCrash, "nocrash", cfg.NoCrash, "do not crash on panics in algorithms")

	flag.Parse()

	if len(nodesRaw) == 0 {
		panic("nodes list cannot be empty")
	}
	cfg.Nodes = strings.Split(nodesRaw, ",")
	if cfg.Self >= len(cfg.Nodes) || cfg.Self < 0 {
		panic(fmt.Sprintf("self (%d) must be in [0,%d)", cfg.Self, len(cfg.Nodes)))
	}
	return NewRuntime(cfg, os.Stderr)
}
