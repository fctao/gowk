package context

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"errors"
	"github.com/mokeoo/gowk/context/query"
	"gopkg.in/go-playground/validator.v9"
	"context"
	"fmt"
)

type Context interface {
	QueryParams(interface{})
	PostForm(interface{})
	RequestJsonBody(interface{})
	GetContext() context.Context
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

type SimpleContext struct {
	request *http.Request
	writer  http.ResponseWriter
	ctx     context.Context
}

func NewContext(w http.ResponseWriter, r *http.Request) (Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(r.Context())
	return &SimpleContext{request: r, writer: w, ctx: ctx}, cancel
}

//query struct , and out struct too
func (ctx *SimpleContext) QueryParams(typ interface{}) {
	q := ctx.request.URL.Query()
	err := query.Unmarshal(q, typ)
	if nil != err {
		panic(err)
	}
	ctx.validate(typ)
}

func (ctx *SimpleContext) PostForm(typ interface{}) {
	form := ctx.request.PostForm
	err := query.Unmarshal(form, typ)
	if nil != err {
		panic(err)
	}
	ctx.validate(typ)
}

func (ctx *SimpleContext) validate(typ interface{}) {
	valid := validator.New()
	err := valid.Struct(typ)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(errors.New("invalid validation params"))
		}
		var invalidMsg string
		for _, err := range err.(validator.ValidationErrors) {
			invalidMsg += fmt.Sprintln(strings.Join([]string{fmt.Sprintf("field:%v", err.Field()),
				fmt.Sprintf("tag:%v", err.ActualTag()),
				fmt.Sprintf("param:%v", err.Param())}, ","))
		}
		panic(errors.New(fmt.Sprintf("validate params has error %s", invalidMsg)))
	}
}

func (ctx *SimpleContext) GetContext() context.Context {
	return ctx.ctx
}

//cast request body to json
func (ctx *SimpleContext) RequestJsonBody(typ interface{}) {
	contentType := ctx.request.Header.Get("Content-Type")
	if strings.Contains(contentType, "json") {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if nil != err {
			panic(err)
		}
		e := json.Unmarshal(body, typ)
		if nil != e {
			panic(e)
		}
	}
	panic(errors.New("request body content-type is not contains json "))
}

func (ctx *SimpleContext) Request() *http.Request {
	return ctx.request
}

func (ctx *SimpleContext) ResponseWriter() http.ResponseWriter {
	return ctx.writer
}
