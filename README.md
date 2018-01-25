### example
1. router.NewContextRouter instance of http.handler
2. run gowk.Server, can use ```kill -USR2 pid``` restart server

### def app
```
func Server(addr string, r *router.ContextRouter, preStartHook func() error) error {
	opt := gracehttp.PreStartProcess(preStartHook)
	return gracehttp.ServeWithOptions(
		[]*http.Server{
			{Addr: addr, Handler: r.Handler()},
		},
		opt,
	)
}
```

### router example
```
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
```


### query params cast to struct and validate params values
#### def struct
```
type Person struct {
	UserId   int    `query:"user_id",json:"user_id",validate:"gte=0,lt=1000"`
	Age      int    `query:"age",json:"age"`
	Password string `query:"password",json:"password"`
}
```
#### def handler
```
	r.ContextHandlerFunc(func(ctx context.Context) router.ResponseEntity {
		p := &Person{}
		ctx.QueryParams(p)
		return &msg.ResponseBody{
			ContentType: msg.Json,
			Data:        p,
		}
	}).Path("/hello")
```


### Performance
```
goos: linux
goarch: amd64
pkg: github.com/mokeoo/gowk
10000	    159824 ns/op
PASS
```

### dependencies
```github.com/facebookgo/grace/gracehttp```

```github.com/gorilla/mux```

```gopkg.in/go-playground/validator.v9```