package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	"github.com/rutkin/gophkeeper/internal/server/core/port"
	mock_port "github.com/rutkin/gophkeeper/internal/server/core/service/mock"
	"github.com/stretchr/testify/require"
)

type BankMatcher struct {
	data bankItem
	arg  domain.BankData
}

func (bm BankMatcher) Matches(x interface{}) bool {
	bank, ok := x.(domain.BankData)
	if !ok {
		return false
	}
	bm.arg = bank
	if bm.data.Cvv != bm.arg.Card.Cvv {
		return false
	}
	if bm.data.Holder != bm.arg.Card.CardHolder {
		return false
	}
	if bm.data.Meta != bm.arg.Ctx.Meta {
		return false
	}
	if bm.data.Number != bm.arg.Card.CardNumber {
		return false
	}
	if bm.data.Title != bm.arg.Ctx.Title {
		return false
	}

	return true
}

func (e BankMatcher) String() string {
	return fmt.Sprintf("matches arg %+v and data %+v", e.arg, e.data)
}

func EqBank(data bankItem) gomock.Matcher {
	return BankMatcher{data: data}
}

func TestKeeper_SetBank(t *testing.T) {
	type fields struct {
		authService   *mock_port.MockAuthService
		keeperService *mock_port.MockKeeper
		tokenService  *mock_port.MockTokenService
	}
	type args struct {
		bank bankItem
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		prepare        func(fields, args)
		expectedStatus int
	}{
		{
			name: "success set bank request",
			args: args{
				bank: bankItem{
					Title:  "title",
					Meta:   "meta",
					Number: "124",
					Holder: "Holder",
					Cvv:    123,
				},
			},
			prepare: func(f fields, a args) {
				f.keeperService.EXPECT().SetBankData(gomock.Any(), EqBank(a.bank))
				f.tokenService.EXPECT().VerifyToken(gomock.Any()).DoAndReturn(
					func(token string) (domain.TokenPayload, error) {
						return domain.TokenPayload{}, nil
					},
				)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthorized set bank request",
			args: args{
				bank: bankItem{},
			},
			prepare: func(f fields, a args) {
				f.tokenService.EXPECT().VerifyToken(gomock.Any()).DoAndReturn(
					func(token string) (domain.TokenPayload, error) {
						return domain.TokenPayload{}, fmt.Errorf("invalid token")
					},
				)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			tt.fields.authService = mock_port.NewMockAuthService(ctrl)
			tt.fields.keeperService = mock_port.NewMockKeeper(ctrl)
			tt.fields.tokenService = mock_port.NewMockTokenService(ctrl)
			handler := NewHandler(tt.fields.authService, tt.fields.keeperService, tt.fields.tokenService)
			tt.prepare(tt.fields, tt.args)

			server := httptest.NewServer(handler)
			body, err := json.Marshal(tt.args.bank)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, server.URL+"/api/keeper/bank", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("authorization", "bearer token")
			req.Header.Set("Content-Type", "application/json")

			client := http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestHandler_GetBank(t *testing.T) {
	type fields struct {
		authService   port.AuthService
		keeperService port.Keeper
		tokenService  port.TokenService
		engine        *gin.Engine
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				authService:   tt.fields.authService,
				keeperService: tt.fields.keeperService,
				tokenService:  tt.fields.tokenService,
				engine:        tt.fields.engine,
			}
			h.GetBank(tt.args.ctx)
		})
	}
}
