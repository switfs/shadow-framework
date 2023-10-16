package logoutmiddleware

import (
	"github.com/astaxie/beego/session"
	"github.com/gin-gonic/gin"
	"github.com/switfs/shadow-framework/middleware/sessionmiddleware"
)

// LogoutFilter is a gin middleware for handle logout request
func LogoutFilter(logoutPath string, globalSessions *session.Manager) gin.HandlerFunc {
	if logoutPath == "" {
		panic("logout path is empty")
	}
	if globalSessions == nil {
		panic("globalSessions is nil")
	}
	return func(c *gin.Context) {
		if c.Request.URL.RequestURI() == logoutPath {
			Log.Info("user logout")
			sess := sessionmiddleware.GetCurrentSession(c)
			sess.Flush()
			globalSessions.SessionDestroy(c.Writer, c.Request)
		}
		c.Next()
	}
}
