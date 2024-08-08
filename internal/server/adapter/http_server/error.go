package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

var errorStatusMap = map[error]int{
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUserExists:                 http.StatusBadRequest,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrBadRequest:                 http.StatusBadRequest,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
}

func validationError(ctx *gin.Context, err error) {
	log.Err(err).Msg("validation error")
	ctx.JSON(http.StatusBadRequest, "validation error")
}

func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	log.Err(err).Msg("response error")
	ctx.JSON(statusCode, err)
}

func handleSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	log.Err(err).Msg("abort response")
	ctx.AbortWithStatusJSON(statusCode, err)
}
