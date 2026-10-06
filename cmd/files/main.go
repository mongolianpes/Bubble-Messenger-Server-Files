package main

import (
	"files/internal/service"
)

func main() {
	filesService := service.NewFilesService()
	go filesService.StartGRPCService()
	go filesService.StartUploader()
	go filesService.StartFilesCleaner()
	select {}
}
