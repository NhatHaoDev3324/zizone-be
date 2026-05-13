package handler

import (
	"net/http"

	"github.com/NhatHaoDev3324/zizone-be/internal/modules/word/dto"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/word/service"
	"github.com/NhatHaoDev3324/zizone-be/pkg/response"
	"github.com/gin-gonic/gin"
)

type wordHandler struct {
	service service.WordService
}

func NewWordHandler(service service.WordService) *wordHandler {
	return &wordHandler{service}
}

func (h *wordHandler) GetListWord(ctx *gin.Context) {
	var params struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Search string `form:"search"`
	}

	if params.Page == 0 {
		params.Page = 1
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	meta, words, err := h.service.GetListWord(params.Page, params.Limit, params.Search)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not get list word: "+err.Error())
		return
	}

	response.SuccessWithMetaAndData(ctx, "Get list word successfully", meta, words)
}

func (h *wordHandler) CreateWord(ctx *gin.Context) {
	var req dto.CreateWordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.service.CreateWord(&req); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not create word: "+err.Error())
		return
	}

	response.SuccessWithData(ctx, "Create word successfully", req)
}

func (h *wordHandler) GetWordByID(ctx *gin.Context) {
	id := ctx.Param("id")
	word, err := h.service.GetWordByID(id)
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "Word not found")
		return
	}

	response.SuccessWithData(ctx, "Get word successfully", word)
}

func (h *wordHandler) UpdateWord(ctx *gin.Context) {
	var req dto.UpdateWordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.service.UpdateWord(&req); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not update word: "+err.Error())
		return
	}

	response.SuccessWithData(ctx, "Update word successfully", req)
}

func (h *wordHandler) DeleteWord(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.DeleteWord(id); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not delete word: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "Delete word successfully")
}
