package main

import (
	"encoding/hex"
	"errors"
)

// splitPieces splits a large binary string of SHA-1 hashes into 20-byte chunks.
func splitPieces(pieces []byte) ([]string, error) {
	const sha1Length = 20

	if len(pieces)%sha1Length != 0 {
		return nil, errors.New("pieces length is not a multiple of 20 bytes")
	}

	var pieceList []string
	for i := 0; i < len(pieces); i += sha1Length {
		piece := pieces[i : i+sha1Length]
		// pieceList = append(pieceList, base64.RawStdEncoding.EncodeToString(piece))
		pieceList = append(pieceList, hex.EncodeToString(piece))

	}

	return pieceList, nil
}
