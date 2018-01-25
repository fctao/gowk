### example
1. router.NewContextRouter instance of http.handler
2. run gowk.Server, can use ```kill -USR2 pid``` restart server

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