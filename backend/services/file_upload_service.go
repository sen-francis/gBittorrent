package services

import (
	"bittorrent/backend/utils"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type FileUploadService struct {
	ctx	context.Context
}

type File struct {
	Data string `json:"file"`
	Err error `json:"error"`
}

var fileUploadService *FileUploadService

func GetFileUploadService() *FileUploadService {
	if fileUploadService == nil {
		fileUploadService = &FileUploadService{} 
	}
	return fileUploadService
}

func (fileUploadService *FileUploadService) Init(ctx context.Context) {
	fileUploadService.ctx = ctx
}

func (fileUploadService *FileUploadService) SelectFile() File {
	filePath, err := runtime.OpenFileDialog(fileUploadService.ctx, runtime.OpenDialogOptions{
		Title: "Select torrent file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Torrent files (*.torrent)",
				Pattern:     "*.torrent",
			}, 
		},	
	})

	if (err != nil) {
		// TODO SEN: blowup 
	}

	utils.ParseTorrentFile(filePath)

	return File{Data: filePath, Err: err}
	
}

