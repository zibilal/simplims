package gingonic

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"simplims/api"
)

// GonicEngine is wrapper type for gin.ApiEngine type
type GonicEngine struct {
	gonicEngine *gin.Engine
	httpServer  *http.Server
}

func NewGonicEngine(address string) *GonicEngine {
	router := new(GonicEngine)
	router.gonicEngine = gin.Default()

	router.httpServer = &http.Server{
		Addr:    address,
		Handler: router.gonicEngine,
	}
	return router
}

func (e *GonicEngine) RegisterVersion(versions ...*api.Version) error {
	for _, version := range versions {
		routeVersion := e.gonicEngine.Group(version.Name())
		for _, r := range version.Router() {
			handlers := e.wrapHandler(r.Handler, r.Middlewares...)
			switch r.Method {
			case http.MethodPost:
				routeVersion.POST(r.Path, handlers...)
			case http.MethodGet:
				routeVersion.GET(r.Path, handlers...)
			case http.MethodPut:
				routeVersion.PUT(r.Path, handlers...)
			case http.MethodDelete:
				routeVersion.DELETE(r.Path, handlers...)
			case http.MethodPatch:
				routeVersion.PATCH(r.Path, handlers...)
			default:
				return errors.New("invalid version " + version.Name() + " unknown method " + r.Method)
			}
		}
	}
	return nil
}

func (e *GonicEngine) wrapHandler(handler api.ApiHandlerFunc, middlewares ...api.ApiHandlerFunc) []gin.HandlerFunc {
	var result []gin.HandlerFunc

	if handler == nil {
		return nil
	}

	result = make([]gin.HandlerFunc, 0)

	for _, m := range middlewares {
		result = append(result, func(c *gin.Context) {
			err := m(WrapGinContext(c))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			}
		})
	}

	result = append(result, func(c *gin.Context) {
		_ = handler(WrapGinContext(c))
	})

	return result
}

func (e *GonicEngine) Execute() error {
	return e.httpServer.ListenAndServe()
}

func (e *GonicEngine) Shutdown(ctx context.Context) error {
	return e.httpServer.Shutdown(ctx)
}

// GonicEngineContext is a wrapper type for gin.Context type
type GonicEngineContext struct {
	ctx *gin.Context
}

func WrapGinContext(ctx *gin.Context) *GonicEngineContext {
	gonicCtx := new(GonicEngineContext)
	gonicCtx.ctx = ctx
	return gonicCtx
}

func (c *GonicEngineContext) BindJSON(output interface{}) error {
	return c.ctx.BindJSON(output)
}

func (c *GonicEngineContext) BindQuery(output interface{}) error {
	return c.ctx.BindQuery(output)
}

func (c *GonicEngineContext) BindUri(output interface{}) error {
	return c.ctx.Bind(output)
}

func (c *GonicEngineContext) BindForm(output interface{}) error {
	return c.ctx.Bind(output)
}

func (c *GonicEngineContext) Set(key string, value interface{}) {
	c.ctx.Set(key, value)
}

func (c *GonicEngineContext) Get(key string) interface{} {
	val, exist := c.ctx.Get(key)
	if !exist {
		return nil
	}

	return val
}

func (c *GonicEngineContext) JSON(status int, response interface{}) error {
	c.ctx.JSON(status, response)
	return nil
}

func (c *GonicEngineContext) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *GonicEngineContext) UnwrapContext() interface{} {
	return c.ctx
}
