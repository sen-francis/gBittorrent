package services

import (
	"bittorrent/backend/torrent"
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type FileUploadService struct {
	ctx	context.Context
}

type FileUploadResponse struct {
	TorrentMetainfo torrent.TorrentMetainfo 
	Err error 
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

func (fileUploadService *FileUploadService) SelectFile() FileUploadResponse {
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
		fmt.Printf("Error occured while uploading file: %s\n", err.Error())
		return FileUploadResponse{ 
			TorrentMetainfo: torrent.TorrentMetainfo{}, 
			Err: err,
		}
	}

	torrentMetainfo, err := torrent.ParseTorrentFile(filePath)
	
	return FileUploadResponse{ 
		TorrentMetainfo: torrentMetainfo, 
		Err: err,
	}

}

