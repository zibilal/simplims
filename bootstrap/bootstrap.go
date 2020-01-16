package bootstrap

import "simplims/appctx"

type Bootstrap interface {
	Init() error
	Run() error
	ApplicationContext() *appctx.AppContext
}
