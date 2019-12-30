package rboot

import (
	"net/http"
	"os"
	"strconv"

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

// 命名路由
func (r *route) Name(name string) *route {
	if r.name != "" {
		logrus.Errorf("route already has name %q, can't set %q", r.name, name)
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

type Router struct {
	mux    *mux.Router
	routes []*route
}

func newRouter() *Router {
	return &Router{mux: mux.NewRouter(), routes: make([]*route, 0)}
}

// 注册新路由
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *route {
	route := &route{path: path, handlerFunc: f}
	r.routes = append(r.routes, route)
	return route
}

// 注册新路由
func (r *Router) Handle(path string, handler http.Handler) *route {
	route := &route{path: path, handler: handler}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) run() {
	// 注册路由
	r.mux.HandleFunc("/", rbootHome)
	for _, route := range r.routes {
		var routeMux *mux.Route
		if route.handler != nil {
			routeMux = r.mux.Handle(route.path, route.handler)
		} else if route.handlerFunc != nil {
			routeMux = r.mux.HandleFunc(route.path, route.handlerFunc)
		} else {
			continue
		}

		if len(route.methods) > 0 {
			routeMux = routeMux.Methods(route.methods...)
		}

		if route.name != "" {
			routeMux = routeMux.Name(route.name)
		}
	}

	// 获取 web 端口
	addr := os.Getenv("WEB_SERVER_ADDR")
	if addr == "" {
		addr = ":7856"
	}

	logrus.Infof("web 服务开启，地址 %s", addr)

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

func rbootHome(w http.ResponseWriter, r *http.Request) {

	var out = `<div style="color: green;width: 100%;text-align: center;margin-top: 10%;font-size: 18px;"><pre style="word-wrap: break-word; white-space: pre-wrap;">` + rbootLogo + `</pre></div>`
	w.Write([]byte(out))
}
