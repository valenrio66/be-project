package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-project/internal/dto"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := GetAuthPayload(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Error: "Unauthorized"})
			c.Abort()
			return
		}

		roleAllowed := false
		for _, role := range allowedRoles {
			if authPayload.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, dto.APIResponse{Error: "Forbidden: You don't have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
