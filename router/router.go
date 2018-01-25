package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/mokeoo/gowk/context"
	"github.com/mokeoo/gowk/msg"
	"github.com/pkg/errors"
	"fmt"
)

type ContextRouter struct {
	*mux.Router
	errorHandler ErrorHandler
}

type ResponseEntity interface {
	Execute(writer http.ResponseWriter) error
}

type StatusCode int

type Handler func(ctx context.Context) ResponseEntity

type ErrorHandler func(err error) ResponseEntity

func NewContextRouter() *ContextRouter {
	muxRouter := mux.NewRouter()
	return &ContextRouter{Router: muxRouter, errorHandler: defaultErrorHandler}
}

func defaultErrorHandler(err error) ResponseEntity {
	return &msg.ResponseBody{ContentType: msg.Json,
		Data: map[string]string{
			"status":  "500",
			"message": err.Error(),
		}}
}

func (r *ContextRouter) ContextHandlerFunc(handler Handler) *mux.Route {
	route := r.NewRoute()
	return route.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.NewContext(writer, request)
		var resp ResponseEntity
		var recoverMsg interface{}
		defer func() {
			recoverMsg = recover()
			if nil != recoverMsg {
				cancel()
			}
			r.handlerRecover(ctx, resp, recoverMsg)
		}()
		resp = handler(ctx)
	})
}

func (r *ContextRouter) ErrorHandlerFunc(handler ErrorHandler) *mux.Router {
	if nil != handler {
		r.errorHandler = handler
	}
	return r.Router
}

func (r *ContextRouter) NotFoundHandlerFunc(handler http.HandlerFunc) *mux.Router {
	if nil != handler {
		r.NotFoundHandler = handler
	}
	return r.Router
}

func (r *ContextRouter) Handler() http.Handler {
	return r.Router
}

func (r *ContextRouter) handlerRecover(ctx context.Context, result ResponseEntity, recoverMsg interface{}) {
	if nil == recoverMsg {
		parse := result.(ResponseEntity)
		err := parse.Execute(ctx.ResponseWriter())
		if nil != err {
			r.errorHandler(err).Execute(ctx.ResponseWriter())
		}
	} else {
		var e error
		switch recoverMsg.(type) {
		case error:
			e = recoverMsg.(error)
		default:
			e = errors.New(fmt.Sprintf("%v", recoverMsg))
		}
		r.errorHandler(e).Execute(ctx.ResponseWriter())
	}
}
