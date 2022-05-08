package dolphin

import (
	"net/http"
	"strings"
	"sync"
)

type RouterConfig struct {
	NotFoundHandler HandlerFunc
}

type Router struct {
	NotFoundHandler HandlerFunc
	nodeTree        map[string]*routerNode
	rm              sync.Mutex
}

type routerNode struct {
	children map[string]*routerNode
	handlers HandlerChain
}

// DefaultNotFoundHandler is the default handler for 404 requests.
func DefaultNotFoundHandler(ctx *Context) {
	ctx.String("Not Found", http.StatusNotFound)
	ctx.Abort()
}

// NewRouter creates and returns a new router.
func NewRouter(config ...RouterConfig) *Router {
	router := &Router{
		NotFoundHandler: DefaultNotFoundHandler,
		rm:              sync.Mutex{},
	}

	if len(config) > 0 {
		cfg := config[0]
		if cfg.NotFoundHandler != nil {
			router.NotFoundHandler = cfg.NotFoundHandler
		}
	}

	return router
}

// Routes returns the handler for this router.
func (router *Router) Routes() HandlerFunc {
	return func(ctx *Context) {
		node := router.getRouterNode(ctx.Method(), ctx.Path())

		if node != nil {
			ctx.Use(node.handlers...)
			ctx.Next()
		} else if router.NotFoundHandler != nil {
			router.NotFoundHandler(ctx)
		}
	}
}

// ANY adds specific path with both GET, POST, PUT, DELETE, HEAD, OPTIONS, and PATCH methods into router.
func (router *Router) ANY(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("DELETE", path, handlers...)
	router.addRouterNode("GET", path, handlers...)
	router.addRouterNode("HEAD", path, handlers...)
	router.addRouterNode("OPTIONS", path, handlers...)
	router.addRouterNode("PATCH", path, handlers...)
	router.addRouterNode("POST", path, handlers...)
	router.addRouterNode("PUT", path, handlers...)

	return router
}

// DELETE adds specific path with DELETE method into router.
func (router *Router) DELETE(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("DELETE", path, handlers...)

	return router
}

// GET adds specific path with GET method into router.
func (router *Router) GET(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("GET", path, handlers...)

	return router
}

// HEAD adds specific path with HEAD method into router.
func (router *Router) HEAD(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("HEAD", path, handlers...)

	return router
}

// OPTIONS adds specific path with OPTIONS method into router.
func (router *Router) OPTIONS(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("OPTIONS", path, handlers...)

	return router
}

// PATCH adds specific path with PATCH method into router.
func (router *Router) PATCH(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("PATCH", path, handlers...)

	return router
}

// POST adds specific path with POST method into router.
func (router *Router) POST(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("POST", path, handlers...)

	return router
}

// PUT adds specific path with PUT method into router.
func (router *Router) PUT(path string, handlers ...HandlerFunc) *Router {
	router.addRouterNode("PUT", path, handlers...)

	return router
}

func (router *Router) addRouterNode(method string, path string, handlers ...HandlerFunc) {
	router.rm.Lock()
	defer router.rm.Unlock()

	tree := router.nodeTree[method]
	if tree == nil {
		tree = new(routerNode)
		tree.children = make(map[string]*routerNode)
		router.nodeTree[method] = tree
	}

	tree.addRouterNode(path, handlers...)
}

func (router *Router) getRouterNode(method string, path string) *routerNode {
	node := router.nodeTree[method]
	if node == nil {
		return nil
	}

	paths := resolvePath(path)
	for _, p := range paths {
		child := node.children[p]

		if child == nil {
			return nil
		}

		node = child
	}

	return node
}

func (root *routerNode) addRouterNode(path string, handlers ...HandlerFunc) {
	node := root
	paths := resolvePath(path)

	for _, p := range paths {
		child := node.children[p]
		if child == nil {
			child = new(routerNode)
			child.children = make(map[string]*routerNode)

			node.children[p] = child
		}

		node = child
	}

	node.handlers = handlers
}

func resolvePath(path string) []string {
	path = strings.Trim(path, "/")

	return strings.Split(path, "/")
}
