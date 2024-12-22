package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/guddu75/bttorrent/p2p"
	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	Infohash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Peices      string `bencode:"pieces"`
	PieceLenght int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

const Port uint16 = 6881

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	bto := bencodeTorrent{}

	err = bencode.Unmarshal(file, &bto)

	if err != nil {
		return TorrentFile{}, err
	}

	return bto.toTorrentFile()

}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, i)

	if err != nil {
		return [20]byte{}, err
	}

	h := sha1.Sum(buf.Bytes())

	return h, nil
}

func (i *bencodeInfo) splitHashes() ([][20]byte, error) {
	hashlen := 20
	buf := []byte(i.Peices)

	if len(buf)%hashlen != 0 {
		err := fmt.Errorf("Recieved malformed pieces of length %d", len(buf))
		return nil, err
	}

	numHashes := len(buf) / hashlen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashlen:(i+1)*hashlen])
	}

	return hashes, nil

}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()

	if err != nil {
		return TorrentFile{}, err
	}

	pieceHashes, err := bto.Info.splitHashes()

	if err != nil {
		return TorrentFile{}, err
	}

	t := TorrentFile{
		Announce:    bto.Announce,
		Infohash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLenght,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return t, nil
}

func (t *TorrentFile) downloadFile(path string) error {
	var peerId [20]byte

	_, err := rand.Read(peerId[:])

	if err != nil {
		return err
	}

	peers, err := t.requestPeers(peerId, Port)

	if err != nil {
		return err
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      peerId,
		InfoHash:    t.Infohash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
	}

}
