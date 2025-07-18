package torrent

import (
	"bittorrent/backend/utils"
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const AZ_CLIENT_PREFIX string = "-GB0001-"

type TorrentState struct {
	Downloaded int64
	Left int64
	Uploaded int64
	PeerId string
	Event string
	DownloadedPieces []int
}

func (torrentMetainfo *TorrentMetainfo) buildTrackerRequest(torrentState *TorrentState) (string, error) {
	trackerRequest, err := url.Parse(torrentMetainfo.Announce)
	if err != nil {
		return "", err
	}

	urlParams := url.Values{
		"info_hash": []string{string(torrentMetainfo.InfoHash[:])},
		"peer_id": []string{string(torrentState.PeerId)},
		"port": []string{"6888"},
		"uploaded": []string{strconv.FormatInt(torrentState.Uploaded, 10)},
		"downloaded": []string{strconv.FormatInt(torrentState.Downloaded, 10)},
		"left": []string{strconv.FormatInt(torrentMetainfo.Size, 10)},
		"compact": []string{"1"},
		"numwant": []string{"50"},
	}

	if torrentState.Event != "" {
		urlParams.Add("event", torrentState.Event)
	}

	trackerRequest.RawQuery = urlParams.Encode()
	return trackerRequest.String(), nil
}

func parseDictionaryModelPeers(peers []map[string]any) ([]Peer, error) {
	var peerList []Peer
	for _, peerValue := range peers {
		peerIdValue, ok := peerValue["peer id"]
		if !ok {
			return peerList, errors.New("peer id missing from TrackerResponse")
		}
		peerId, ok := peerIdValue.(string)
		if !ok {
			return peerList, errors.New("peer id key in TrackerResponse was not a string")
		}
		
		ipValue, ok := peerValue["ip"]
		if !ok {
			return peerList, errors.New("ip missing from TrackerResponse")
		}
		ip, ok := ipValue.(string)
		if !ok {
			return peerList, errors.New("ip key in TrackerResponse was not a string")
		}
		ipAddress := net.IP([]byte(ip))

		portValue, ok := peerValue["ip"]
		if !ok {
			return peerList, errors.New("port missing from TrackerResponse")
		}
		port, ok := portValue.(int64)
		if !ok {
			return peerList, errors.New("port key in TrackerResponse was not an int")
		}

		peer := Peer {PeerId: peerId, IpAddress: ipAddress, Port: uint(port)}
		peerList = append(peerList, peer)
	}

	return peerList, nil
}

func parseBinaryModelPeers(peers string) ([]Peer, error) {
	byteArr := []byte(peers)
	var peerList []Peer
	if len(byteArr) < 6 {
		return []Peer{}, errors.New("Binary model peers has unexpected format")
	}
	for len(byteArr) >= 6 {
		peer := byteArr[:6]
		ipAddress := net.IP(peer[:4])
		port := uint(binary.BigEndian.Uint16(peer[4:]))
		peerList = append(peerList, Peer{ IpAddress: ipAddress, Port: port} )
		byteArr = byteArr[6:]
	
	}
	return peerList, nil
}

func parseTrackerResponse(dictionary map[string]any) (int64, []Peer, error) {
	peers, ok := dictionary["peers"]
	if !ok {
		return 0, []Peer{}, errors.New("No peers key found in tracker response.") 
	}

	var peerList []Peer
	peersDict, err := castAnyToSliceOfMap(peers)
	if err != nil {
		if peersStr, ok := peers.(string); !ok {	
			return 0, []Peer{}, errors.New("peers value in tracker response was not in expected binary or dictionary model.") 
		} else {
			peerList, err = parseBinaryModelPeers(peersStr)
			if err != nil {
				return 0, []Peer{}, err	
			}
		}
	} else {
		peerList, err = parseDictionaryModelPeers(peersDict)
		if err != nil {
			return 0, []Peer{}, err	
		}
	}

	intervalAny, ok := dictionary["interval"]
	if !ok {
		return 0, []Peer{}, errors.New("No interval key found in tracker response.")  
	}

	interval, ok :=  intervalAny.(int64); 
	if !ok {
		return 0, []Peer{}, errors.New("interval value in tracker response is not of type int")
	}

	return interval, peerList, nil
}

func (torrentMetainfo *TorrentMetainfo) fetchPeers(peerCh chan []Peer, torrentStateCh chan TorrentState) (error) {	
	for {
		fmt.Println("Fetching new peers")
		torrentState := <-torrentStateCh

		trackerRequestUrl, err := torrentMetainfo.buildTrackerRequest(&torrentState)
		if err != nil {
			return err
		}

		resp, err := http.Get(trackerRequestUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
			
		if resp.StatusCode != 200 {
			errorString := fmt.Sprintf("Tracker request returned non-OK response. StatusCode=%d", resp.StatusCode)
			return errors.New(errorString)
		}
		//We Read the response body on the line below.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err 
		}
		//Convert the body to type string
		sb := string(body)
		reader := bufio.NewReader(strings.NewReader(sb))
		decoded, err :=	utils.Decode(reader)
		if err != nil {
			return err
		}

		dictionary, ok := decoded.(map[string]any)
		if !ok { 
			return errors.New("TrackerResponse is invalid.")
		}

		interval, peersList, err := parseTrackerResponse(dictionary)
		if err != nil {
			return err	
		}
		fmt.Printf("Found %d potential peer(s)\n", len(peersList))

		peerCh <- peersList
		time.Sleep(time.Duration(interval) * time.Second)
	}
}


func (torrentMetainfo *TorrentMetainfo) BuildScrapeRequest() (string, error) {
	splitAnnounce := strings.Split(torrentMetainfo.Announce,"/")
	text := splitAnnounce[len(splitAnnounce) - 1]
	if !strings.HasPrefix(text, "announce") {
		return "", errors.New("Tracker has no scrape convention")
	}

	scrapeUrl := strings.Replace(text, "announce", "scrape", 1)
	splitAnnounce[len(splitAnnounce) - 1] = scrapeUrl
	scrapeUrl = strings.Join(splitAnnounce, "/")
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
