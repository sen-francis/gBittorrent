package main

import (
	"bittorrent/backend/services"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fileUploadService := services.GetFileUploadService()
	trackerService := services.GetTrackerService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "bittorrent",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			fileUploadService.Init(ctx)
			trackerService.Init(ctx)
		},
		Bind: []interface{}{
			fileUploadService,
			trackerService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
