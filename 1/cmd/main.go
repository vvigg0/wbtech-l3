package main

import (
	"log"

	"github.com/vvigg0/wbtech-l3/l3/1/cmd/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
