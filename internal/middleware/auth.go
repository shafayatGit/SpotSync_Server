package middleware

import (
	"net/http"
	"strings"

	"spotsync/internal/auth"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates JWT Bearer tokens in Request Headers or Cookies and populates context variables.
func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenStr string

			// 1. Try to get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenStr = parts[1]
				}
			}

			// 2. If not found in header, try to get from "token" cookie
			if tokenStr == "" {
				cookie, err := c.Cookie("token")
				if err == nil {
					tokenStr = cookie.Value
				}
			}

			// 3. Fallback: try "access_token" cookie
			if tokenStr == "" {
				cookie, err := c.Cookie("access_token")
				if err == nil {
					tokenStr = cookie.Value
				}
			}

			// If still empty, reject request
			if tokenStr == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": "Authorization token is required (in header or cookie)",
				})
			}

			claims, err := jwtService.ValidateToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": "Invalid or expired token",
					"errors":  err.Error(),
				})
			}

			// Store user information in context
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

// RoleMiddleware checks if the user's role is in the list of allowed roles.
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleVal := c.Get("role")
			role, ok := roleVal.(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"message": "Access forbidden: role not found",
				})
			}

			for _, r := range allowedRoles {
				if r == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"message": "Access forbidden: insufficient permissions",
			})
		}
	}
}
