package main

import (
	"files/internal/service"
)

func main() {
	filesService := service.NewFilesService()
	go filesService.StartSaver()
	go filesService.StartUploader()
	go filesService.StartFilesCleaner()
}
