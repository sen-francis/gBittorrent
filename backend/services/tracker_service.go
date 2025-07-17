package services

import (
	"bittorrent/backend/torrent"
	"bittorrent/backend/utils"
	"bufio"
	"context"
	"io"
	"net/http"
	"strings"
	"errors"
)

type TrackerService struct {
	ctx context.Context
}

type TrackerResponse struct {
	Interval int64
	Peers string
	Err error
}

type TrackerScrapeResponse struct {
	Downloaded int32
	Seeders int32
	Leechers int32
	Name string	
	Err error
}

var trackerService *TrackerService

func GetTrackerService() *TrackerService {
	if trackerService == nil {
		trackerService = &TrackerService{} 
	}
	return trackerService
}

func (trackerService *TrackerService) Init(ctx context.Context) {
	trackerService.ctx = ctx
}



func parseTrackerScrapeResponse(response string, infoHash string) TrackerScrapeResponse {
	reader := bufio.NewReaderSize(strings.NewReader(response), len(response))

	result, err := utils.Decode(reader)
	if err != nil {
		return TrackerScrapeResponse{Err: err}
	}

	dictionary, ok := result.(map[string]any)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}

	if _, ok := dictionary["files"]; !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}

	files, ok := dictionary["files"].(map[string]any)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}

	if len(files) != 1 {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}

	if _, ok := files[infoHash]; !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}

	torrentData, ok := files[infoHash].(map[string]any)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}

	var trackerScrapeResponse TrackerScrapeResponse;
	seeders, ok := torrentData["complete"].(int64)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}
	trackerScrapeResponse.Seeders = int32(seeders)

	leechers, ok := torrentData["incomplete"].(int64)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}
	trackerScrapeResponse.Leechers = int32(leechers)

	downloaded, ok := torrentData["downloaded"].(int64)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}
	trackerScrapeResponse.Downloaded = int32(downloaded)
	if torrentData["name"] != nil {
		trackerScrapeResponse.Name = torrentData["name"].(string)
	}

	return trackerScrapeResponse 
}

func (trackerService *TrackerService) ScrapeTracker(torrentMetainfo *torrent.TorrentMetainfo) TrackerScrapeResponse {
	scrapeRequestUrl, err := torrentMetainfo.BuildScrapeRequest()
	if err != nil {
		return TrackerScrapeResponse{ Err: err}
	}

	resp, err := http.Get(scrapeRequestUrl)
	if err != nil {
		return TrackerScrapeResponse{ Err: err}
	}
	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackerScrapeResponse{Err: err}
	}
	// todo
	//Convert the body to type string
	sb := string(body)



	trackerScrapeResponse := parseTrackerScrapeResponse(sb, string(torrentMetainfo.InfoHash[:]))

	return trackerScrapeResponse
}
