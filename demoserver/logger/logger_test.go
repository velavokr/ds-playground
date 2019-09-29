package logger

import (
	"bytes"
	"github.com/velavokr/gdaf/demoserver/utils"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	type param struct {
		verbose, verboseOnly bool
	}
	for _, p := range []param{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	} {
		type fn struct {
			name  string
			apply func(l *Logger, b *bytes.Buffer) bool
		}
		for _, f := range []fn{
			{"Print", func(l *Logger, b *bytes.Buffer) bool {
				l.Println(p.verboseOnly, "hello", "world", []byte("a"))
				return utils.ContainsAll(b.String(), "greet", "logger_test.go:", `hello world {slice=[97],hex=61,ascii="a"}`, "\n")
			}},
			{"Printf", func(l *Logger, b *bytes.Buffer) bool {
				l.Printf(p.verboseOnly, "%s %s", "hello", "world")
				return utils.ContainsAll(b.String(), "greet", "logger_test.go:", `hello world`, "\n") && !strings.Contains(b.String(), "%s")
			}},
			{"Output", func(l *Logger, b *bytes.Buffer) bool {
				l.Output(l.Caller(0), p.verboseOnly, "hello", "world", []byte("a"))
				return utils.ContainsAll(b.String(), "greet", "logger_test.go:", `hello world {slice=[97],hex=61,ascii="a"}`, "\n")
			}},
			{"Outputf", func(l *Logger, b *bytes.Buffer) bool {
				l.Outputf(l.Caller(0), p.verboseOnly, "%s %s", "hello", "world")
				return utils.ContainsAll(b.String(), "greet", "logger_test.go:", `hello world`, "\n") && !strings.Contains(b.String(), "%s")
			}},
		} {
			t.Run(f.name, func(t *testing.T) {
				b := bytes.Buffer{}
				l := NewLogger(p.verbose, &b, "greet")
				c := f.apply(l, &b)
				if !c && (p.verbose || !p.verboseOnly) || c && (!p.verbose && p.verboseOnly) {
					t.Errorf("%#v %s", p, b.String())
				}
			})
		}
	}
}
