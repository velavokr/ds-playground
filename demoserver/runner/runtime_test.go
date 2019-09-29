package runner

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
)

func TestRuntime_RunSync(t *testing.T) {
	p := "catch me if you can"
	b := bytes.Buffer{}
	e := NewRuntime(UserCfg{}, &b)
	e.RunGuarded(func() {
		panic(p)
	}, 0, "first")
	e.RunGuarded(func() {
	}, 0, "second")
	e.Run(func() {
		panic(p)
	}, 0, "third")
	for _, s := range []string{"[node 0] ", " runtime_test.go:"} {
		if strings.Count(b.String(), s) != 6 {
			t.Error(b.String(), s)
		}
	}
	if strings.Count(b.String(), " guarded ") != 4 {
		t.Error(b.String())
	}
	if strings.Count(b.String(), " run ") != 3 {
		t.Error(b.String())
	}
	for _, s := range []string{" panic ", " first", " second", " third"} {
		if strings.Count(b.String(), s) != 2 {
			t.Error(b.String(), s)
		}
	}
	if strings.Count(b.String(), " done ") != 1 {
		t.Error(b.String(), )
	}
}

func TestRuntime_RunAsync(t *testing.T) {
	p := "catch me if you can"
	b := bytes.Buffer{}
	e := NewRuntime(UserCfg{}, &b)
	wg := sync.WaitGroup{}
	wg.Add(1)
	e.RunAsync(func(ctx context.Context) {
		defer wg.Done()
		select {
		case <-ctx.Done():
		}
	}, 0, "first")
	wg.Add(1)
	e.RunAsync(func(ctx context.Context) {
		defer wg.Done()
		panic(p)
	}, 0, "second")
	e.Cancel()
	e.WaitAll()
	wg.Wait()
	if strings.Count(b.String(), p) != 1 {
		t.Error()
	}
	for _, s := range []string{"[node 0] ", " async ", " runtime_test.go:"} {
		if strings.Count(b.String(), s) != 4 {
			t.Error(b.String())
		}
	}
	for _, s := range []string{" run ", " first", " second"} {
		if strings.Count(b.String(), s) != 2 {
			t.Error(b.String(), s)
		}
	}
	for _, s := range []string{" panic ", " done "} {
		if strings.Count(b.String(), s) != 1 {
			t.Error(b.String(), s)
		}
	}
}
