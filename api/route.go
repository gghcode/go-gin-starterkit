package api

import "github.com/gin-gonic/gin"

type HandlerFunc = gin.HandlerFunc
type RouteInfos []RouteInfo

type RouteInfo struct {
	Method string
	Path   string
	Handle HandlerFunc
}

func Route(method, path string, handle HandlerFunc) RouteInfo {
	return RouteInfo{
		Method: method,
		Path:   path,
		Handle: handle,
	}
}
