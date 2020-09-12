package rboot

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type route struct {
	// 路由名称
	name string

	// 访问路径
	path string

	// 路由访问方式
	methods []string

	// 处理函数
	handlerFunc func(http.ResponseWriter, *http.Request)
	handler     http.Handler
}

// Name 为命名路由
func (r *route) Name(name string) *route {
	if r.name != "" {
		logrus.WithFields(logrus.Fields{
			"mod": `rboot`,
		}).Errorf("route already has name %q, can't set %q", r.name, name)
	} else {
		r.name = name
	}

	return r
}

// 设置 methods
func (r *route) Methods(methods ...string) *route {
	r.methods = methods
	return r
}

// Router 包含了路由处理器 mux 和已经注册的所有路由集合，支持中间件
type Router struct {
	mux    *mux.Router
	routes []*route

	// 中间件
	middlewares []func(http.Handler) http.Handler
}

// newRouter 创建一个路由实例
func newRouter() *Router {
	return &Router{
		mux:         mux.NewRouter(),
		routes:      make([]*route, 0),
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Use 注册中间件，和 *mux.Router.Use 用法相同
func (r *Router) Use(mwf ...func(http.Handler) http.Handler) {
	for _, fn := range mwf {
		r.middlewares = append(r.middlewares, fn)
	}
}

// HandleFunc 为路径 path 注册一个新的路由处理函数
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *route {
	ro := &route{path: path, handlerFunc: f}
	r.routes = append(r.routes, ro)
	return ro
}

// Handle 为路径 path 注册一个新路由
func (r *Router) Handle(path string, handler http.Handler) *route {
	ro := &route{path: path, handler: handler}
	r.routes = append(r.routes, ro)
	return ro
}

func (r *Router) run() {
	// 注册路由
	r.mux.HandleFunc("/", webHome)
	r.mux.HandleFunc("/ipv4", remoteIPV4)

	for _, ro := range r.routes {
		var routeMux *mux.Route
		if ro.handler != nil {
			routeMux = r.mux.Handle(ro.path, ro.handler)
		} else if ro.handlerFunc != nil {
			routeMux = r.mux.HandleFunc(ro.path, ro.handlerFunc)
		} else {
			continue
		}

		if len(ro.methods) > 0 {
			routeMux = routeMux.Methods(ro.methods...)
		}

		if ro.name != "" {
			routeMux = routeMux.Name(ro.name)
		}
	}

	if len(r.middlewares) > 0 {
		for _, middleware := range r.middlewares {
			r.mux.Use(mux.MiddlewareFunc(middleware))
		}
	}

	r.mux.StrictSlash(true)

	// 获取 web 端口
	port := os.Getenv("WEB_SERVER_PORT")
	if port == "" {
		port = "5689"
	}

	var addr = ":" + port

	fmt.Println("web 服务开启，地址 ", addr)

	isTls, _ := strconv.ParseBool(os.Getenv("WEB_SERVER_TLS"))
	if isTls {
		cert := os.Getenv("WEB_SERVER_CERT")
		certKey := os.Getenv("WEB_SERVER_CERT_KEY")
		if err := http.ListenAndServeTLS(addr, cert, certKey, r.mux); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(addr, r.mux); err != nil {
			panic(err)
		}
	}
}

func remoteIPV4(w http.ResponseWriter, r *http.Request) {
	var remoteAddr string
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0]); ip != "" {
		remoteAddr = ip
	}

	if ip := strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" {
		remoteAddr = ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		remoteAddr = ip
	}

	w.Write([]byte(remoteAddr))
}

func webHome(w http.ResponseWriter, r *http.Request) {

	var out = `<div style="color: green;width: 100%;text-align: center;margin-top: 10%;font-size: 18px;"><pre style="word-wrap: break-word; white-space: pre-wrap;">` + rbootLogo + `</pre></div>`
	w.Write([]byte(out))
}
