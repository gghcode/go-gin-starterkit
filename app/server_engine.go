package app

type ServerEngine interface {
	RegisterRoute(route RouteInfo)
	Run(addr string) error
}
