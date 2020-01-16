package v3

import (
	"net/http"
	"simplims/api"
	"simplims/api/v3/handler"
)

func VersionOneApi() *api.Version {
	endpoints := []api.Endpoint{
		{
			Path:    "/ping",
			Method:  http.MethodGet,
			Handler: handler.PingApi,
			Middlewares: []api.ApiHandlerFunc{
				handler.MiddlewareTest,
			},
		},
	}
	return api.NewVersion("v3", endpoints)
}
