package demoserver

import (
	"context"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/demoserver/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ApiHandler interface {
	HandleApiCall(url *url.URL, body []byte) (responseBody []byte, err error)
}

func RunServer(rt *runner.Runtime, handler ApiHandler) {
	StartServer(rt, handler)
	rt.WaitAll()
}

func StartServer(rt *runner.Runtime, handler ApiHandler) {
	server := &http.Server{
		Addr: rt.Cfg.Http,
		Handler: &httpHandler{
			rt:      rt,
			handler: handler,
		},
	}

	rt.RunAsync(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			if err := server.Close(); err != nil {
				panic(err)
			}
		}
	}, runner.VerboseLog, "server stopper")

	rt.RunAsync(func(ctx context.Context) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}, runner.VerboseLog, "http server", rt.Cfg.Http)
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rt.Run(func() {
		body := []byte{}

		switch r.Method {
		case "POST", "PUT":
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			if b != nil {
				body = b
			}
		}

		var res []byte
		var err error
		h.rt.RunGuarded(func() {
			res, err = h.handler.HandleApiCall(r.URL, body)
		}, runner.VerboseLog|runner.ExitOnPanic, "handle request", r.URL, "with body", body)

		if err != nil {
			w.WriteHeader(500)
			if err := utils.WriteAll(w, []byte(err.Error())); err != nil {
				panic(err)
			}
			panic(err)
		}

		w.WriteHeader(200)
		if err := utils.WriteAll(w, res); err != nil {
			panic(err)
		}
	}, runner.VerboseLog, "serve request", r.Method, r.URL)
}

type httpHandler struct {
	rt      *runner.Runtime
	handler ApiHandler
}
