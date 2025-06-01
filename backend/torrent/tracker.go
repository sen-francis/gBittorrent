package torrent

import (
	"errors"
	"math/rand/v2"
	"net/url"
	"strings"
)

const AZ_CLIENT_PREFIX string = "-GB0001-"

type TrackerResponse struct {
	interval int64
	peers string
}

func generatePeerId() string {
	random := ""
	for i := 0; i < 12; i++ {
		random += string(rand.IntN(10))
	}
	return AZ_CLIENT_PREFIX + random
}

func (torrentInfo *TorrentInfo) getTotalTorrentBytes() int64 {
	totalBytes := int64(0)
	for _, fileInfo := range torrentInfo.FileInfoList {
		totalBytes += fileInfo.Length
	}

	return totalBytes
}

func (torrentMetainfo *TorrentMetainfo) BuildTrackerRequest() (string, error) {
	trackerRequest, err := url.Parse(torrentMetainfo.Announce)
	if err != nil {
		return "", err
	}

	urlParams := url.Values{
		"info_hash": []string{string(torrentMetainfo.InfoHash[:])},
		"peer_id": []string{generatePeerId()},
		"port": []string{"6888"},
		"uploaded": []string{"0"},
		"downloaded": []string{"0"},
		"left": []string{string(torrentMetainfo.Info.getTotalTorrentBytes())},
		"compact": []string{"1"},
		"event": []string{"started"},
	}

	trackerRequest.RawQuery = urlParams.Encode()
	return trackerRequest.String(), nil
}

func (torrentMetainfo *TorrentMetainfo) BuildScrapeRequest() (string, error) {
	splitAnnounce := strings.Split(torrentMetainfo.Announce,"\\")
	text := splitAnnounce[len(splitAnnounce) - 1]
	if strings.HasPrefix(text, "announce") {
		scrapeUrl := strings.Replace(text, "announce", "scrape", 1)
		splitAnnounce[len(splitAnnounce) - 1] = scrapeUrl
		scrapeUrl = strings.Join(splitAnnounce, "\\")
		scrapeRequest, err := url.Parse(scrapeUrl)
		if err != nil {
			return "", err
		}

		urlParams := url.Values{
			"info_hash": []string{string(torrentMetainfo.InfoHash[:])},
		}

		scrapeRequest.RawQuery = urlParams.Encode()
		return scrapeRequest.String(), nil
	}

	return "", errors.New("Tracker has no scrape convention")
}
