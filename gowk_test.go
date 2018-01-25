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
	"encoding/json"
	"errors"
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

type Person struct {
	UserId   int    `query:"user_id",json:"user_id",validate:"gte=0,lt=1000"`
	Age      int    `query:"age",json:"age"`
	Password string `query:"password",json:"password"`
}

func TestQueryParams(t *testing.T) {
	r := router.NewContextRouter()

	r.ContextHandlerFunc(func(ctx context.Context) router.ResponseEntity {
		p := &Person{}
		ctx.QueryParams(p)
		return &msg.ResponseBody{
			ContentType: msg.Json,
			Data:        p,
		}
	}).Path("/hello")
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/hello?user_id=110&age=18&password=admin")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	rp := &Person{}
	e := json.Unmarshal(body, rp)
	if nil != e {
		t.Error(e)
	}

	if rp.UserId != 110 || rp.Age != 18 || rp.Password != "admin" {
		t.Error("response error")
	}
}

func TestPanicError(t *testing.T) {
	r := router.NewContextRouter()
	r.ContextHandlerFunc(func(ctx context.Context) router.ResponseEntity {
		panic(errors.New("error test"))
	}).Path("/error")

	ts := httptest.NewServer(r)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/error")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	t.Log(string(body))
}

func BenchmarkPanicError(b *testing.B) {
	r := router.NewContextRouter()
	r.ContextHandlerFunc(func(ctx context.Context) router.ResponseEntity {
		panic(errors.New("error test"))
	}).Path("/error")

	ts := httptest.NewServer(r)
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/error")
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if !strings.Contains(string(body), "error test") {
			b.Error("has error")
		}
	}
}

func BenchmarkServer(b *testing.B) {
	ts := httptest.NewServer(getTestContextRouter())
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/hello")
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if !strings.Contains(string(body), "hello world") {
			b.Error("has error")
		}

	}
}
