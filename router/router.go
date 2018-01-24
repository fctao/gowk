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

type HandlerResultParse interface {
	Execute(writer http.ResponseWriter) error
}

type StatusCode int

type Handler func(ctx context.Context) HandlerResultParse

type ErrorHandler func(err error) HandlerResultParse

func NewContextRouter() *ContextRouter {
	muxRouter := mux.NewRouter()
	return &ContextRouter{Router: muxRouter,}
}

func (r *ContextRouter) ViewHandlerFunc(handler Handler) *mux.Route {
	route := r.NewRoute()
	return route.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.NewContext(writer, request)
		result := handler(ctx)
		recoverMsg := recover()
		if nil != recoverMsg {
			cancel()
		}
		r.handlerRecover(ctx, result, recoverMsg)
	})
}

func (r *ContextRouter) ErrorHandlerFunc(handler ErrorHandler) *mux.Router {
	if nil != handler {
		r.errorHandler = handler
	} else {
		r.errorHandler = func(err error) HandlerResultParse {
			return &msg.ResponseBody{ContentType: msg.Json,
				Data: map[string]string{
					"status":  "500",
					"message": err.Error(),
				}}
		}
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

func (r *ContextRouter) handlerRecover(ctx context.Context, result HandlerResultParse, recoverMsg interface{}) {
	if nil == recoverMsg {
		parse := result.(HandlerResultParse)
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