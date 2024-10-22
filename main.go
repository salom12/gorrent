package main

import (
	"os"

	"github.com/salom12/gorrent/bencode"
)

type TorrentMeta struct {
	Comment      string `json:"comment" bencode:"comment"`
	CreatedBy    string `json:"created_by" bencode:"created by"`
	CreationDate int64  `json:"creation_date" bencode:"creation date"` // timestamps (UNIX time)

	Info struct {
		Name        string `json:"name" bencode:"name"`
		Length      int64  `json:"length" bencode:"length"`
		PieceLength int64  `json:"piece_length" bencode:"piece length"`
		Pieces      []byte `json:"pieces" bencode:"pieces"` // Binary data, use []byte
	} `json:"info" bencode:"info"` // Info is typically a dictionary with various fields

	UrlList []string `json:"url_list" bencode:"url-list"`
}

func main() {
	// read input torrent
	torrentFile, err := os.ReadFile("test_data/archlinux-2024.10.01-x86_64.iso.torrent")
	if err != nil {
		panic(err)
	}

	// parse it to struct
	meta := TorrentMeta{}
	bencode.Unmarshal(torrentFile, &meta)

	// convert it back to bencode
	data, err := bencode.Marshal(meta)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("output.torrent", data, 0644); err != nil {
		panic(err)
	}
}
