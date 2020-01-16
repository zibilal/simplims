package echo

import (
	"context"
	"errors"
	"github.com/labstack/echo"
	logger "github.com/zibilal/logwrapper"
	"github.com/zibilal/simpleapi/api"
	"net/http"
	"reflect"
	"strconv"
)

// EchoEngine is wrapper type for echo.Echo type
type EchoEngine struct {
	address    string
	echoEngine *echo.Echo
}

func NewEchoEngine(address string) *EchoEngine {
	router := new(EchoEngine)
	router.echoEngine = echo.New()
	router.address = address

	return router
}

func (e *EchoEngine) RegisterVersions(versions ...*api.Version) error {

	for _, version := range versions {
		routeVersion := e.echoEngine.Group(version.Name())
		for _, r := range version.Router() {
			switch r.Method {
			case http.MethodPost:
				middles := e.wrapMiddlewares(r.Middlewares)
				routeVersion.POST(r.Path, echo.HandlerFunc(func(c echo.Context) error {
					return r.Handler(WrapEchoEngineContext(c))
				}), middles...)
			case http.MethodGet:
				middles := e.wrapMiddlewares(r.Middlewares)
				routeVersion.GET(r.Path, echo.HandlerFunc(func(c echo.Context) error {
					return r.Handler(WrapEchoEngineContext(c))
				}), middles...)
			case http.MethodPut:
				middles := e.wrapMiddlewares(r.Middlewares)
				routeVersion.PUT(r.Path, echo.HandlerFunc(func(c echo.Context) error {
					return r.Handler(WrapEchoEngineContext(c))
				}), middles...)
			case http.MethodDelete:
				middles := e.wrapMiddlewares(r.Middlewares)
				routeVersion.DELETE(r.Path, echo.HandlerFunc(func(c echo.Context) error {
					return r.Handler(WrapEchoEngineContext(c))
				}), middles...)
			case http.MethodPatch:
				middles := e.wrapMiddlewares(r.Middlewares)
				routeVersion.PATCH(r.Path, echo.HandlerFunc(func(c echo.Context) error {
					return r.Handler(WrapEchoEngineContext(c))
				}), middles...)
			default:
				return errors.New("invalid version " + version.Name() + " unknown method " + r.Method)
			}
		}
	}

	return nil
}

func (e *EchoEngine) wrapMiddlewares(middlewares []api.ApiHandlerFunc) []echo.MiddlewareFunc {
	if middlewares == nil || len(middlewares) == 0 {
		return nil
	}

	result := make([]echo.MiddlewareFunc, len(middlewares))
	for i, m := range middlewares {
		result[i] = func(next echo.HandlerFunc) echo.HandlerFunc {

			return func(c echo.Context) error {
				if err := m(WrapEchoEngineContext(c)); err != nil {
					return err
				}
				return next(c)
			}
		}
	}

	return result
}

func (e *EchoEngine) Execute() error {
	return e.echoEngine.Start(e.address)
}

func (e *EchoEngine) Shutdown(ctx context.Context) error {
	return e.echoEngine.Shutdown(ctx)
}

type EchoEngineContext struct {
	ctx echo.Context
}

func WrapEchoEngineContext(ctx echo.Context) *EchoEngineContext {
	echoCtx := new(EchoEngineContext)
	echoCtx.ctx = ctx

	return echoCtx
}

func (c *EchoEngineContext) BindJSON(output interface{}) error {
	return c.ctx.Bind(output)
}

func (c *EchoEngineContext) BindQuery(output interface{}) error {
	return c.simpleIterateType("query", output)
}

func (c *EchoEngineContext) BindUri(output interface{}) error {
	return c.simpleIterateType("uri", output)
}

func (c *EchoEngineContext) BindForm(output interface{}) error {
	return c.simpleIterateType("form", output)
}

func (c *EchoEngineContext) UnwrapContext() interface{} {
	return c.ctx
}

func (c *EchoEngineContext) Set(key string, value interface{}) {
	c.ctx.Set(key, value)
}

func (c *EchoEngineContext) Get(key string) interface{} {
	return c.ctx.Get(key)
}

func (c *EchoEngineContext) SetStatusCode(status int) {
	c.ctx.Response().WriteHeader(status)
}

func (c *EchoEngineContext) JSON(code int, response interface{}) error {
	return c.ctx.JSON(code, response)
}

func (c *EchoEngineContext) SetHeader(key, value string) {
	c.ctx.Response().Header().Set(key, value)
}

func (c *EchoEngineContext) simpleIterateType(tagName string, output interface{}) error {
	ival := reflect.ValueOf(output)
	ityp := reflect.TypeOf(output)

	var tempStr string

	if ival.Kind() == reflect.Struct {
		for i := 0; i < ival.NumField(); i++ {
			fval := ival.Field(i)

			if fval.IsValid() {
				continue
			}

			ftyp := ityp.Field(i)
			tag := ftyp.Tag.Get(tagName)

			switch tagName {
			case "form":
				tempStr = c.ctx.FormValue(tag)
			case "uri":
				tempStr = c.ctx.Param(tag)
			case "query":
				tempStr = c.ctx.QueryParam(tag)
			case "json":
				return errors.New("[EchoEngineContext] for json tag please use the default ctx.Bind method")
			default:
				return errors.New("[EchoEngineContext] unknown tag type " + tagName)
			}

			if tempStr != "" {
				// check the type of the field
				switch fval.Kind() {
				case reflect.String, reflect.Interface:
					fval.Set(reflect.ValueOf(tempStr))
				case reflect.Int:
					tmp, err := strconv.ParseInt(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iint := int(tmp)
					fval.Set(reflect.ValueOf(iint))
				case reflect.Int8:
					tmp, err := strconv.ParseInt(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iint8 := int8(tmp)
					fval.Set(reflect.ValueOf(iint8))
				case reflect.Int16:
					tmp, err := strconv.ParseInt(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iint16 := int16(tmp)
					fval.Set(reflect.ValueOf(iint16))
				case reflect.Int32:
					tmp, err := strconv.ParseInt(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iint32 := int32(tmp)
					fval.Set(reflect.ValueOf(iint32))
				case reflect.Int64:
					tmp, err := strconv.ParseInt(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iint64 := int64(tmp)
					fval.Set(reflect.ValueOf(iint64))
				case reflect.Uint:
					tmp, err := strconv.ParseUint(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iuint := uint(tmp)
					fval.Set(reflect.ValueOf(iuint))
				case reflect.Uint8:
					tmp, err := strconv.ParseUint(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iuint8 := uint8(tmp)
					fval.Set(reflect.ValueOf(iuint8))
				case reflect.Uint16:
					tmp, err := strconv.ParseUint(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iuint16 := uint16(tmp)
					fval.Set(reflect.ValueOf(iuint16))
				case reflect.Uint32:
					tmp, err := strconv.ParseUint(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iuint32 := uint32(tmp)
					fval.Set(reflect.ValueOf(iuint32))
				case reflect.Uint64:
					tmp, err := strconv.ParseUint(tempStr, 10, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					iuint64 := uint64(tmp)
					fval.Set(reflect.ValueOf(iuint64))
				case reflect.Float32:
					tmp, err := strconv.ParseFloat(tempStr, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					ifloat32 := float32(tmp)
					fval.Set(reflect.ValueOf(ifloat32))
				case reflect.Float64:
					tmp, err := strconv.ParseFloat(tempStr, 64)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					ifloat64 := float64(tmp)
					fval.Set(reflect.ValueOf(ifloat64))
				case reflect.Bool:
					tmp, err := strconv.ParseBool(tempStr)
					if err != nil {
						logger.Error("[EchoEngineContext] ", err.Error())
					}
					ibool := bool(tmp)
					fval.Set(reflect.ValueOf(ibool))
				default:
					return errors.New("[EchoEngineContext] unsupported type " + ftyp.Type.String())
				}
			}
		}
	} else {
		return errors.New("[EchoEngineContext] only accept output of type struct")
	}

	return nil
}
