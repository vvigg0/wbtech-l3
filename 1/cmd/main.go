package main

import (
	"github.com/vvigg0/wbtech-l3/l3/1/cmd/app"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	if err := app.Run(); err != nil {
		zlog.Logger.Fatal().Msg(err.Error())
	}
}
