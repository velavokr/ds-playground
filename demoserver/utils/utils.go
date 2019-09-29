package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func ToDeadline(tout time.Duration) time.Time {
	return time.Now().Add(tout)
}

func WriteAll(w io.Writer, b []byte) error {
	for len(b) > 0 {
		n, err := w.Write(b)
		if err != nil {
			return err
		}
		b = b[n:]
	}
	return nil
}

func Quote(b []byte) string {
	if b == nil {
		return "{nil}"
	}
	return fmt.Sprintf("{slice=%v,hex=%s,ascii=%s}", b, hex.EncodeToString(b), strconv.QuoteToASCII(string(b)))
}

func Sprint(cs ...interface{}) string {
	b := bytes.Buffer{}
	for i, c := range cs {
		if i > 0 {
			b.WriteByte(' ')
		}
		bb, ok := c.([]byte)
		if ok {
			b.WriteString(Quote(bb))
			continue
		}
		b.WriteString(fmt.Sprint(c))
	}
	return b.String()
}

func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

func Less(b [][]byte) func(i, j int) bool {
	return func(i, j int) bool {
		return bytes.Compare(b[i], b[j]) < 0
	}
}

func RandomFreePort() int {
	for i := 0; i < 32768; i++ {
		p := 32768+rand.Intn(28231)
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err == nil {
			_ = l.Close()
			return p
		}
	}
	return -1
}
