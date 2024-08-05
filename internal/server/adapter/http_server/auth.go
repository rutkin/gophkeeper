package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

func (h *Handler) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.User{
		Name:     domain.UserName(req.Name),
		Password: req.Password,
	}

	err := h.authService.Register(ctx, user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}

type loginRequest struct {
	Name     string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

type loginResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

func (h *Handler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	token, err := h.authService.Login(ctx, domain.User{Name: domain.UserName(req.Name), Password: req.Password})
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAuthResponse(token)

	handleSuccess(ctx, rsp)
}

func newAuthResponse(token domain.Token) loginResponse {
	return loginResponse{string(token)}
}
