package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) UpdateNotificationStatus(ctx *ginext.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		zlog.Logger.Error().Msgf("invalid id: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": err.Error()})
		return
	}

	err = h.srvc.UpdateNotificationStatus(id, "cancelled")
	if err != nil {
		zlog.Logger.Error().Msgf("не удалось обновить статус уведомления: %s", err.Error())
		if errors.Is(err, repository.ErrNoNotification) {
			ctx.JSON(http.StatusBadRequest, ginext.H{
				"err": err.Error(),
			})
		}
		ctx.JSON(http.StatusInternalServerError, ginext.H{"err": err.Error()})
		return
	}

	zlog.Logger.Info().Msgf("DELETE запрос по уведомлению с id=%v прошел успешно", id)
	ctx.JSON(http.StatusOK, ginext.H{"res": "успешно удалено"})
}
