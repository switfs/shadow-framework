package logoutmiddleware

import (
	"github.com/astaxie/beego/session"
	"github.com/switfs/shadow-framework/logger"
	"github.com/switfs/shadow-framework/middleware"
)

var (
	Log            *logger.Logger
	globalSessions *session.Manager
)

const (
	LOGOUT         = "logout"
	LOGOUT_HANDLER = "LogoutHandler"
)

func init() {
	Log = logger.InitLog()
	Log.Info("DefaultLogoutUrlRegistry init")
	middleware.RegisterMiddlewareHandler(LOGOUT_HANDLER, newDefaultLogoutHandler)
}
