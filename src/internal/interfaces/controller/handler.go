package controller

import api "platform-go/generated/api"

// Handler combines generated user endpoints and authentication endpoints.
type Handler struct {
	*UserController
	*AuthController
}

func NewHandler(users *UserController, auth *AuthController) *Handler {
	return &Handler{UserController: users, AuthController: auth}
}

var _ api.ServerInterface = (*Handler)(nil)
