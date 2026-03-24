package handler

import (
	"net/http"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) CreateNotification(ctx *ginext.Context) {
	c := ctx.Request.Context()
	var req dto.CreateNotificationsRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": "invalid JSON"})
		zlog.Logger.Error().Msgf("ошибка парсинга JSON: %v", err)
		return
	}

	ids, err := h.srvc.CreateNotifications(c, req)
	if err != nil {
		if len(ids) == 0 {
			ctx.JSON(http.StatusBadRequest, ginext.H{"err": err})
			return
		}
		ctx.JSON(http.StatusOK, ginext.H{"res": ids, "err": err})
		return
	}

	zlog.Logger.Info().Msg("POST запрос прошел успешно")
	ctx.JSON(http.StatusOK, ginext.H{"res": ids})
}
