package gowk

import (
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
	"github.com/mokeoo/gowk/router"
)
//addr is listen address
//router.ContextRouter is instance of handler
//support preStartHook, use "kill -USR2 1104" kill current pid and start new pid , first run preStartHost fun
func Server(addr string, r *router.ContextRouter, preStartHook func() error) error {
	opt := gracehttp.PreStartProcess(preStartHook)
	return gracehttp.ServeWithOptions(
		[]*http.Server{
			{Addr: addr, Handler: r.Handler()},
		},
		opt,
	)
}
