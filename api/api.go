package api

import "context"

type ApiEngine interface {
	RegisterVersion(versions ...*Version) error
	Execute() error
	Shutdown(ctx context.Context) error
}

type Endpoint struct {
	Path        string
	Method      string
	Handler     ApiHandlerFunc
	Middlewares []ApiHandlerFunc
}

type Version struct {
	endpoints   []Endpoint
	versionName string
}

func NewVersion(versionName string, endpoints []Endpoint) *Version {
	v := new(Version)
	v.endpoints = endpoints
	v.versionName = versionName
	return v
}

func (v *Version) Name() string {
	return v.versionName
}

func (v *Version) AddEndpoint(path, method string, handler ApiHandlerFunc, middlewares ...ApiHandlerFunc) {
	v.endpoints = append(v.endpoints, Endpoint{
		Path:        path,
		Method:      method,
		Handler:     handler,
		Middlewares: middlewares,
	})
}

type ApiHandlerFunc func(EngineContext) error

type ApiResponder interface {
	Response() interface{}
}

func (v *Version) Router() []Endpoint {
	return v.endpoints
}

type EngineContext interface {
	BindJSON(output interface{}) error
	BindQuery(output interface{}) error
	BindUri(output interface{}) error
	BindForm(output interface{}) error
	UnwrapContext() interface{}
	Set(key string, value interface{})
	Get(key string) interface{}
	SetHeader(key, value string)
	JSON(code int, value interface{}) error
}
