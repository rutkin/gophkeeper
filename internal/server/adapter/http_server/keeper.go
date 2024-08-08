package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

type itemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type listItemsResponse struct {
	Items []itemResponse `json:"items"`
}

func (h *Handler) ListItems(ctx *gin.Context) {
	payload := getAuthPayload(ctx)
	meta, err := h.keeperService.ListAll(ctx, payload.ID)
	if err != nil {
		if err == domain.ErrNotFound {
			handleSuccess(ctx, nil)
			return
		}
		handleError(ctx, err)
		return
	}

	var resp listItemsResponse
	for _, m := range meta {
		resp.Items = append(resp.Items, itemResponse{ID: string(m.ID), Name: m.Title, Type: string(m.Type)})
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) UploadFile(ctx *gin.Context) {
	reader, err := ctx.Request.MultipartReader()
	if err != nil {
		log.Err(err).Msg("failed read multipart request")
		handleError(ctx, err)
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FormName() == "file" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, part)
			if err != nil {
				log.Err(err).Msg("failed to file buffer")
				handleError(ctx, err)
				return
			}

			payload := getAuthPayload(ctx)
			dataCtx := domain.DataContext{
				ID:     domain.DataID(uuid.NewString()),
				UserID: payload.ID,
				Type:   domain.BinaryType,
				Title:  part.FileName(),
			}
			err = h.keeperService.SetBinaryData(ctx, domain.BinaryData{Ctx: dataCtx, Data: buf.Bytes()})
			if err != nil {
				log.Err(err).Msg("failed to set binary data")
				handleError(ctx, err)
				return
			}
			handleSuccess(ctx, "")
			return
		}
	}
	handleError(ctx, domain.ErrInvalidToken)
}

func (h *Handler) DownloadFile(ctx *gin.Context) {
	dataID := ctx.Param("id")
	payload := getAuthPayload(ctx)
	data, err := h.keeperService.GetBinaryData(ctx, domain.DataContext{ID: domain.DataID(dataID), UserID: payload.ID})
	if err != nil {
		log.Err(err).Msg("failed to get binary data")
		handleError(ctx, err)
		return
	}
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Length", strconv.Itoa(len(data.Data)))
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", data.Ctx.Title))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Data(http.StatusOK, "application/octet-stream", data.Data)
}

type credentialsItem struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Title    string `json:"title"`
	Meta     string `json:"meta"`
}

func (h *Handler) SetCredentials(ctx *gin.Context) {
	var req credentialsItem
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Err(err).Msg("failed to bind cred request")
		handleError(ctx, err)
		return
	}

	payload := getAuthPayload(ctx)

	err = h.keeperService.SetCredentialsData(ctx, domain.CredentialsData{
		Ctx: domain.DataContext{
			ID:     domain.DataID(uuid.NewString()),
			UserID: payload.ID,
			Meta:   req.Meta,
			Title:  req.Title,
			Type:   domain.CredentialsType,
		},
		Cred: domain.Credentials{
			Username: req.Name,
			Password: req.Password,
		},
	})

	if err != nil {
		log.Err(err).Msg("failed to set credentials")
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "")
}

func (h *Handler) GetCredentials(ctx *gin.Context) {
	id := ctx.Param("id")
	payload := getAuthPayload(ctx)
	data, err := h.keeperService.GetCredentialsData(ctx, domain.DataContext{ID: domain.DataID(id), UserID: payload.ID})
	if err != nil {
		log.Err(err).Msg("failed to get credentials")
		handleError(ctx, err)
		return
	}
	resp := credentialsItem{
		Name:     data.Cred.Username,
		Password: data.Cred.Password,
		Title:    data.Ctx.Title,
		Meta:     data.Ctx.Meta,
	}
	handleSuccess(ctx, resp)
}

type bankItem struct {
	Title  string
	Meta   string
	Number string
	Holder string
	Cvv    int
}

func (h *Handler) SetBank(ctx *gin.Context) {
	var req bankItem
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Err(err).Msg("failed to get bank request")
		handleError(ctx, err)
		return
	}
	payload := getAuthPayload(ctx)
	err = h.keeperService.SetBankData(ctx, domain.BankData{
		Ctx: domain.DataContext{
			ID:     domain.DataID(uuid.NewString()),
			UserID: payload.ID,
			Meta:   req.Meta,
			Title:  req.Title,
			Type:   domain.BankType,
		},
		Card: domain.Card{
			CardNumber: req.Number,
			CardHolder: req.Holder,
			Cvv:        req.Cvv,
		},
	})
	if err != nil {
		log.Err(err).Msg("failed to set bank data")
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}

func (h *Handler) GetBank(ctx *gin.Context) {
	id := ctx.Param("id")
	payload := getAuthPayload(ctx)
	data, err := h.keeperService.GetBankData(ctx, domain.DataContext{ID: domain.DataID(id), UserID: payload.ID})
	if err != nil {
		log.Err(err).Msg("failed to get bank data")
		handleError(ctx, err)
		return
	}
	resp := bankItem{
		Number: data.Card.CardNumber,
		Holder: data.Card.CardHolder,
		Cvv:    data.Card.Cvv,
		Title:  data.Ctx.Title,
		Meta:   data.Ctx.Meta,
	}
	handleSuccess(ctx, resp)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	payload := getAuthPayload(ctx)
	err := h.keeperService.Delete(ctx, domain.DataContext{ID: domain.DataID(id), UserID: payload.ID})
	if err != nil {
		log.Err(err).Msg("failed to delete item")
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, nil)
}
