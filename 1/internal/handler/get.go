package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) GetNotificationStatus(ctx *ginext.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка обработки id: %v", err)
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": "ID должен быть числом"})
		return
	}

	status, err := h.srvc.GetNotificationStatus(id)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка получения статуса уведомления: %v", err)
		if errors.Is(err, repository.ErrNoNotification) {
			ctx.JSON(http.StatusBadRequest, ginext.H{
				"err": err,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ginext.H{
			"err": "внутренняя ошибка сервера",
		})
		return
	}

	zlog.Logger.Info().Msgf("GET запрос по уведомлению с id=%v прошел успешно", id)
	ctx.JSON(http.StatusOK, ginext.H{"res": status})
}
