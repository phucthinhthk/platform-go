package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"

	"platform-go/internal/common/ctxkeys"
	"platform-go/internal/pkg/types"
	"platform-go/internal/usecases"
)

func JWTAuth(auth *usecases.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/refresh" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}
		token, err := c.Cookie("jwt")
		if err != nil {
			unauthorized(c)
			return
		}
		claims, err := auth.ParseAccess(token)
		if err != nil {
			unauthorized(c)
			return
		}
		account, err := auth.CurrentUser(c.Request.Context(), claims.Subject)
		if err != nil {
			unauthorized(c)
			return
		}
		c.Set(ctxkeys.LoginAdministratorKey, types.LoginAdministrator{ID: uint(account.ID), Name: account.Name, Email: account.Email, AccountType: account.AccountType})
		c.Next()
	}
}

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Every authenticated API must present the CSRF token returned by login.
		// OPTIONS is left open so browser CORS preflight requests can complete.
		if c.Request.Method == http.MethodOptions || c.Request.URL.Path == "/login" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}
		cookie, err := c.Cookie("csrf")
		header := c.GetHeader("X-CSRF-Token")
		if err != nil || cookie == "" || header == "" || subtle.ConstantTimeCompare([]byte(cookie), []byte(header)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token is invalid"})
			return
		}
		c.Next()
	}
}
func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}
