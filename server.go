package main

import (
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
	"github.com/mokeoo/gowk/router"
	"github.com/mokeoo/gowk/context"
	"github.com/mokeoo/gowk/msg"
)

func main() {
	r := router.NewContextRouter()
	r.ViewHandlerFunc(func(ctx context.Context) router.HandlerResultParse {
		return &msg.ResponseBody{ContentType: msg.Json, Data: map[string]string{
			"status_code": "001",
			"msg":         "hello world",
		}}
	}).Path("/")

	gracehttp.Serve(
		&http.Server{Addr: ":8080", Handler: r.Handler()},
	)
}
