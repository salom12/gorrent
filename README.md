
# Gorrent

This is a work-in-progress implementation of a Torrent client and Bencode parser in Go. The purpose of this project is to deepen my understanding of the Torrent protocol and the Bencode format.

## Status
- [x] bencode
- [ ] torrent client

## Installation

To use this package in your project, you can download it via `go get`:

```bash
go get github.com/salom12/gorrent
```

## Usage
Here’s a simple example of how to use  bencode package for encoding and decoding Bencode data (check tests for parsing examples):

```go
package main

import (
	"os"

	"github.com/salom12/gorrent/bencode"
)

type TorrentMeta struct {
	Comment      string `bencode:"comment"`
	CreatedBy    string `bencode:"created by"`
	CreationDate int64  `bencode:"creation date"` // timestamps (UNIX time)

	Info struct {
		Name        string `bencode:"name"`
		Length      int64  `bencode:"length"`
		PieceLength int64  `bencode:"piece length"`
		Pieces      []byte `bencode:"pieces"` // Binary data, use []byte
	} `bencode:"info"` // Info is typically a dictionary with various fields

	UrlList []string `bencode:"url-list"`
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

```

## Running Tests
```bash
go test ./...
```

## License

This project is licensed under the MIT License.