package torrent

type TorrentState struct {
	Downloaded int
	Left int
	Uploaded int
}

type Peer struct {
	PeerId int
	IpAddress string
	Port int
}

func (torrentMetainfo *TorrentMetainfo) startDownload() {
	torrentMetainfo.BuildTrackerRequest()
}

func fetchPeers() {


}
