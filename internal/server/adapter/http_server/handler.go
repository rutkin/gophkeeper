package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rutkin/gophkeeper/internal/server/core/port"
)

type Handler struct {
	authService   port.AuthService
	keeperService port.Keeper
	tokenService  port.TokenService
	engine        *gin.Engine
}

func NewHandler(authService port.AuthService, keeperService port.Keeper, tokenService port.TokenService) *Handler {
	engine := gin.New()

	handler := &Handler{authService: authService, keeperService: keeperService, tokenService: tokenService, engine: engine}
	handler.init()

	return handler
}

func (h *Handler) init() {
	h.engine.Use(gin.Recovery())
	h.engine.Use(gin.Logger())

	h.engine.POST("api/register", h.Register)
	h.engine.POST("api/login", h.Login)

	keeper := h.engine.Group("api/keeper").Use(authMiddleware(h.tokenService))
	{
		keeper.GET("/", h.ListItems)
		keeper.POST("/file", h.UploadFile)
		keeper.GET("/file/:id", h.DownloadFile)
		keeper.POST("/credentials", h.SetCredentials)
		keeper.GET("/credentials/:id", h.GetCredentials)
		keeper.POST("/bank", h.SetBank)
		keeper.GET("/bank/:id", h.GetBank)
		keeper.POST("/delete/:id", h.Delete)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.engine.ServeHTTP(w, r)
}
