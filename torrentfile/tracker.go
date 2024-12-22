package torrentfile

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/guddu75/bttorrent/peers"
	"github.com/jackpal/bencode-go"
)

type bencodeTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildTrackerURL(peerId [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)

	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.Infohash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()

	return base.String(), nil

}

func (t *TorrentFile) requestPeers(peerId [20]byte, port uint16) ([]peers.Peer, error) {
	url, err := t.buildTrackerURL(peerId, port)

	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}

	resp, err := c.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	trackerResponse := bencodeTrackerResponse{}

	err = bencode.Unmarshal(resp.Body, &trackerResponse)

	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(trackerResponse.Peers))

}
