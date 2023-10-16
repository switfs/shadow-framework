package securitymiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/switfs/shadow-framework/middleware/sessionmiddleware"
	"github.com/switfs/shadow-framework/security"
)

//SecurityTokenResolver 从session中解析出authentication token，用于认证当前的请求上下文
func SecurityTokenResolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从session中取出token, 放入request上下文中
		sess := sessionmiddleware.GetCurrentSession(c)
		authentication := sess.Get(security.SHADOW_SECURITY_TOKEN)

		if authentication != nil {
			c.Set(security.SHADOW_SECURITY_TOKEN, authentication)
		}

		c.Next()

		//请求返回前将token再写回session中
		if auth, exist := c.Get(security.SHADOW_SECURITY_TOKEN); exist {
			sess.Set(security.SHADOW_SECURITY_TOKEN, auth)
		}

	}
}
