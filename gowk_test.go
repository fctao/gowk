package gowk

import (
	"testing"
	"github.com/mokeoo/gowk/router"
	"github.com/mokeoo/gowk/context"
	"github.com/mokeoo/gowk/msg"
	"net/http"
	"io/ioutil"
	"strings"
	"net/http/httptest"
)

func getTestContextRouter() *router.ContextRouter {
	r := router.NewContextRouter()
	r.ContextHandlerFunc(func(ctx context.Context) router.ResponseEntity {
		return &msg.ResponseBody{
			ContentType: msg.Json,
			Data: map[string]string{
				"status":  "200",
				"message": "hello world",
			},
		}
	}).Path("/hello")
	return r
}

func BenchmarkServer(b *testing.B) {
	ts := httptest.NewServer(getTestContextRouter())
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL+"/hello")
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if !strings.Contains(string(body), "hello world") {
			b.Error("has error")
		}

	}
}
