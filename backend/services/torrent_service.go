package services

import (
	"context"
	"bittorrent/backend/torrent"
)

type TorrentService struct {
	ctx context.Context
}


var torrentService *TorrentService

func GetTorrentService() *TorrentService {
	if torrentService == nil {
		torrentService = &TorrentService{}
	}
	return torrentService
}

func (torrentService *TorrentService) Init(ctx context.Context) {
	torrentService.ctx = ctx
}

func (torrentService *TorrentService) DownloadTorrent(torrentMetainfo torrent.TorrentMetainfo) {
	torrentMetainfo.StartDownload(torrentService.ctx)
}

