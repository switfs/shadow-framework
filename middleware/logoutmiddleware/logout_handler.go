package logoutmiddleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/switfs/shadow-framework/middleware"
)

//TDefaultLogoutHandler logout middlewareHandler implementation
type TDefaultLogoutHandler struct{}

func newDefaultLogoutHandler() middleware.IMiddlewareHandler {
	return new(TDefaultLogoutHandler)
}

//Handle handle redirect
func (handler *TDefaultLogoutHandler) Handle(c *gin.Context) {
	Log.Debugln("logoutHandler redirect to /")
	c.Redirect(http.StatusFound, "/")
}
