package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"platform-go/internal/pkg/types"
	"platform-go/internal/usecases"
)

type AuthController struct {
	usecase       *usecases.AuthUsecase
	secureCookies bool
}

func NewAuthController(usecase *usecases.AuthUsecase, secureCookies bool) *AuthController {
	return &AuthController{usecase: usecase, secureCookies: secureCookies}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var body loginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	_, access, refresh, csrf, err := ac.usecase.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	ac.setCookies(c, access, refresh, csrf)
	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"csrfToken": csrf},
		"message": "Login successful",
	})
}

func (ac *AuthController) Refresh(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		ac.clearCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	access, nextRefresh, csrf, err := ac.usecase.Refresh(c.Request.Context(), refresh)
	if err != nil {
		ac.clearCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	ac.setCookies(c, access, nextRefresh, csrf)
	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"csrfToken": csrf},
		"message": "Token refreshed",
	})
}

func (ac *AuthController) Logout(c *gin.Context) { ac.clearCookies(c); c.Status(http.StatusNoContent) }

func (ac *AuthController) Me(c *gin.Context) {
	value, exists := c.Get("login_administrator")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	account, ok := value.(types.LoginAdministrator)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid authenticated context"})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (ac *AuthController) setCookies(c *gin.Context, access, refresh, csrf string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", access, 24*60*60, "/", "", ac.secureCookies, true)
	c.SetCookie("refresh_token", refresh, 7*24*60*60, "/refresh", "", ac.secureCookies, true)
	// The client must read this non-HttpOnly value and mirror it in X-CSRF-Token.
	c.SetCookie("csrf", csrf, 24*60*60, "/", "", ac.secureCookies, false)
}
func (ac *AuthController) clearCookies(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", ac.secureCookies, true)
	c.SetCookie("refresh_token", "", -1, "/refresh", "", ac.secureCookies, true)
	c.SetCookie("csrf", "", -1, "/", "", ac.secureCookies, false)
}
