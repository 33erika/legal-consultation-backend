package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix       = "Bearer "
	UserContextKey    = "user"
)

func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			utils.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			utils.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		user, err := authSvc.ValidateToken(token)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, utils.ErrAuthTokenExpired, "Token已过期或无效")
			c.Abort()
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}

// GetCurrentUser 获取当前登录用户
func GetCurrentUser(c *gin.Context) *service.UserInfo {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	return user.(*service.UserInfo)
}

// RequireRole 角色权限中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "无权限访问此页面")
		c.Abort()
	}
}
