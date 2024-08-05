package httpserver

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	"github.com/rutkin/gophkeeper/internal/server/core/port"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationType      = "bearer"
)

func authMiddleware(ts port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := domain.ErrEmptyAuthorizationHeader
			handleAbort(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := domain.ErrInvalidAuthorizationHeader
			handleAbort(ctx, err)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != authorizationType {
			err := domain.ErrInvalidAuthorizationType
			handleAbort(ctx, err)
			return
		}

		accessToken := fields[1]
		payload, err := ts.VerifyToken(accessToken)
		if err != nil {
			log.Err(err).Msg("failed to verify token")
			handleAbort(ctx, domain.ErrInvalidToken)
			return
		}

		setAuthPayload(ctx, payload)
		ctx.Next()
	}
}
