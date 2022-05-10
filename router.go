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
	handlers        HandlerChain
	nodeTree        map[string]*routerNode
	rm              sync.Mutex
}

type routerNode struct {
	children      map[string]*routerNode
	handlers      HandlerChain
	wildcardChild *routerNode
	pathVarName   string
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
		handlers:        make(HandlerChain, 0),
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

func (router *Router) Use(handler ...HandlerFunc) *Router {
	if len(handler) > 0 {
		router.handlers = append(router.handlers, handler...)
	}

	return router
}

// Routes returns the handler for this router.
func (router *Router) Routes() HandlerFunc {
	return func(ctx *Context) {
		pathVariables := make(map[string]string)
		node := router.getRouterNode(ctx.Method(), ctx.Path(), pathVariables)

		if node != nil {
			if len(router.handlers) > 0 {
				ctx.Use(router.handlers...)
			}
			if len(pathVariables) > 0 {
				for k, v := range pathVariables {
					ctx.addPathVariable(k, v)
				}
			}

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

func (router *Router) getRouterNode(method string, path string, pathVariables map[string]string) *routerNode {
	node := router.nodeTree[method]
	if node == nil {
		return nil
	}

	paths := resolvePath(path)
	for _, p := range paths {
		child := node.children[p]

		if child == nil {
			if node.wildcardChild == nil {
				return nil
			}

			child = node.wildcardChild
			pathVariables[child.pathVarName] = p
		}

		node = child
	}

	return node
}

func (root *routerNode) addRouterNode(path string, handler ...HandlerFunc) {
	var child *routerNode
	node := root
	paths := resolvePath(path)

	for _, p := range paths {
		isWildcard := false

		if strings.HasPrefix(p, ":") {
			child = node.wildcardChild
			isWildcard = true
		} else {
			child = node.children[p]
		}

		if child == nil {
			child = new(routerNode)
			child.children = make(map[string]*routerNode)

			if isWildcard {
				node.wildcardChild = child
				node.pathVarName = p[1:]
			} else {
				node.children[p] = child
			}
		}

		node = child
	}

	node.handlers = make(HandlerChain, 0, len(handler))
	if len(handler) > 0 {
		node.handlers = append(node.handlers, handler...)
	}
}

func resolvePath(path string) []string {
	path = strings.Trim(path, "/")

	return strings.Split(path, "/")
}
