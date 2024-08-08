package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

var authorizationPayloadKey = "payload"

func setAuthPayload(ctx *gin.Context, token domain.TokenPayload) {
	ctx.Set(authorizationPayloadKey, token)
}

func getAuthPayload(ctx *gin.Context) domain.TokenPayload {
	return ctx.MustGet(authorizationPayloadKey).(domain.TokenPayload)
}
