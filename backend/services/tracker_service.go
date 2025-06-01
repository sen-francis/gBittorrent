package services

import (
	"bittorrent/backend/torrent"
	"bittorrent/backend/utils"
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"errors"
)

type TrackerService struct {
	ctx context.Context
}

type TrackerResponse struct {
	interval int64
	peers string
	Err error
}

type TrackerScrapeResponse struct {
	downloaded int64
	seeders int64
	leechers int64
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

func (trackerService *TrackerService) FetchPeers(torrentMetainfo *torrent.TorrentMetainfo) (TrackerResponse, error) {
	trackerRequestUrl, err := torrentMetainfo.BuildTrackerRequest()
	if err != nil {
		return TrackerResponse{}, err
	}

	resp, err := http.Get(trackerRequestUrl)
	if err != nil {
		return TrackerResponse{}, err
	}
	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackerResponse{}, nil
	}
	//Convert the body to type string
	sb := string(body)
	
	trackerResponse := TrackerResponse{}

	fmt.Printf(sb)
	return trackerResponse, nil
}

func parseTrackerScrapeResponse(response string) TrackerScrapeResponse {
	reader := bufio.NewReaderSize(strings.NewReader(response), len(response))

	result, err := utils.Decode(reader)
	if err != nil {
		return TrackerScrapeResponse{Err: err}
	}

	dictionary, ok := result.(map[string]any)
	if !ok {
		return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
	}
	files, ok := dictionary["files"];

	if ok {
		if _, ok := files.([]any); !ok {
			return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
		} else {
			return  TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
		}

	}

	if files, ok := dictionary["files"]; ok {
		if _, ok := files.([]any); !ok {
			return TrackerScrapeResponse{Err: errors.New("Torrent scrape response invalid")}
		}

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

	

	trackerScrapeResponse := TrackerScrapeResponse{ }

	fmt.Printf(sb)
	return trackerScrapeResponse
}
