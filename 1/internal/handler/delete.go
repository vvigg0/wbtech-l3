package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) DeleteNotification(ctx *ginext.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		zlog.Logger.Error().Msgf("invalid id: %v", err)
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": err})
		return
	}

	err = h.srvc.DeleteNotification(id)

	if err != nil {
		zlog.Logger.Error().Msgf("не удалось удалить уведомление: %v", err)
		if errors.Is(err, repository.ErrNoNotification) {
			ctx.JSON(http.StatusNotFound, ginext.H{
				"err": err,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ginext.H{"err": err})
		return
	}

	zlog.Logger.Info().Msgf("DELETE запрос по уведомлению с id=%v прошел успешно", id)
	ctx.JSON(http.StatusOK, ginext.H{"res": "успешно удалено"})
}
