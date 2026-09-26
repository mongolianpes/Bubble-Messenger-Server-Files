package main

import (
	"files/internal/service"
)

func main() {
	filesService := service.NewFilesService()
	go filesService.StartDownloader()
	go filesService.StartUploader()
	go filesService.StartFilesCleaner()
}
