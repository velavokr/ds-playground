package demoserver

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/velavokr/gdaf/demoserver/runner"
	"github.com/velavokr/gdaf/demoserver/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestRunServer(t *testing.T) {
	b := bytes.Buffer{}
	p := utils.RandomFreePort()
	func() {
		rt := runner.NewRuntime(runner.UserCfg{Http: fmt.Sprintf("127.0.0.1:%d", p)}, &b)
		h := &handler{}
		StartServer(rt, h)

		resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/xxx", p), "plain/text", bytes.NewBufferString("Hello"))
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Error(resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "This is fine" {
			t.Error(string(body))
		}
		if err := resp.Body.Close(); err != nil {
			t.Error(err)
		}
		rt.Cancel()
		rt.WaitAll()

		if len(h.items) != 1 {
			t.Error(h.items)
		}
		if h.items[0].url.Path != "/xxx" {
			t.Error(h.items[0].url)
		}
		if string(h.items[0].body) != "Hello" {
			t.Error(h.items[0].body)
		}
	}()

	func() {
		rt := runner.NewRuntime(runner.UserCfg{Http: fmt.Sprintf("127.0.0.1:%d", p)}, &b)
		h := &handler{err: errors.New("nope")}
		StartServer(rt, h)

		resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/xxx", p), "plain/text", bytes.NewBufferString("Hello"))
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 500 {
			t.Error(resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) != "nope" {
			t.Error(string(body))
		}
		if err := resp.Body.Close(); err != nil {
			t.Error(err)
		}
		rt.Cancel()
		rt.WaitAll()

		if len(h.items) != 1 {
			t.Error(h.items)
		}
		if h.items[0].url.Path != "/xxx" {
			t.Error(h.items[0].url)
		}
		if string(h.items[0].body) != "Hello" {
			t.Error(h.items[0].body)
		}
	}()
}

type item struct {
	url  url.URL
	body []byte
}

type handler struct {
	items []item
	err   error
}

func (h *handler) HandleApiCall(url *url.URL, body []byte) (responseBody []byte, err error) {
	h.items = append(h.items, item{url: *url, body: append([]byte{}, body...)})
	if h.err != nil {
		return nil, h.err
	}
	return []byte("This is fine"), nil
}
