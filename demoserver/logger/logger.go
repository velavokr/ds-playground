package logger

import (
	"fmt"
	"github.com/velavokr/gdaf/demoserver/utils"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

func NewLogger(verbose bool, out io.Writer, prefix string) *Logger {
	if len(prefix) > 0 && prefix[len(prefix)-1] != ' ' {
		prefix += " "
	}
	return &Logger{
		out:     out,
		prefix:  prefix,
		verbose: verbose,
	}
}

func (d *Logger) Println(verboseOnly bool, v ...interface{}) {
	if d.verbose || !verboseOnly {
		d.Output(d.Caller(1), verboseOnly, v...)
	}
}

func (d *Logger) Printf(verboseOnly bool, format string, v ...interface{}) {
	if d.verbose || !verboseOnly {
		d.Outputf(d.Caller(1), verboseOnly, format, v...)
	}
}

func (d *Logger) Outputf(caller string, verboseOnly bool, format string, v ...interface{}) {
	if d.verbose || !verboseOnly {
		b := []byte{}
		b = append(b, []byte(d.logPrefix(caller))...)
		b = append(b, []byte(fmt.Sprintf(format, v...))...)
		b = append(b, '\n')
		d.doOutput(b)
		_ = utils.WriteAll(d.out, b)
	}
}

func (d *Logger) Output(caller string, verboseOnly bool, v ...interface{}) {
	if d.verbose || !verboseOnly {
		b := []byte{}
		b = append(b, []byte(d.logPrefix(caller))...)
		b = append(b, []byte(utils.Sprint(v...))...)
		b = append(b, '\n')
		d.doOutput(b)
	}
}

func (d *Logger) Caller(depth int) string {
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return fmt.Sprintf("%s:%d", file, line)
}

func (d *Logger) logPrefix(caller string) string {
	t := time.Now()
	return fmt.Sprintf(
		"%s %02d:%02d:%02d.%09d %s: ", d.prefix, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), caller)
}

func (d *Logger) doOutput(b []byte) {
	if _, ok := d.out.(*os.File); ok {
		_ = utils.WriteAll(d.out, b)
	} else {
		d.mtx.Lock()
		defer d.mtx.Unlock()
		_ = utils.WriteAll(d.out, b)
	}
}

type Logger struct {
	mtx     sync.Mutex
	out     io.Writer
	prefix  string
	verbose bool
}
