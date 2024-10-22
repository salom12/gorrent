package bencode

import (
	"errors"
	"fmt"
	"strconv"
)

type Bencode struct {
	input  []byte
	offset int
	Result []any
}

// parseString parses a bencoded string in the form "length:string"
func (b *Bencode) parseString() (any, error) {
	colonIndex := -1
	for i := b.offset; i < len(b.input); i++ {
		if b.input[i] == ':' {
			colonIndex = i
			break
		}
	}
	if colonIndex == -1 {
		return nil, errors.New("invalid string format, missing ':'")
	}

	// Use strconv.ParseInt for parsing length directly from byte slice
	length, err := strconv.ParseInt(string(b.input[b.offset:colonIndex]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid string length: %v", err)
	}

	start := colonIndex + 1
	end := start + int(length)
	if end > len(b.input) {
		return nil, errors.New("string length exceeds input size")
	}

	byteValue := b.input[start:end]
	b.offset = end

	return byteValue, nil
}

// parseInt parses a bencoded integer in the form "i<number>e"
func (b *Bencode) parseInt() (any, error) {
	b.offset++ // Skip 'i'
	endIndex := -1
	for i := b.offset; i < len(b.input); i++ {
		if b.input[i] == 'e' {
			endIndex = i
			break
		}
	}
	if endIndex == -1 {
		return nil, errors.New("invalid integer format, missing 'e'")
	}

	// Use strconv.ParseInt to parse the integer directly from the byte slice
	number, err := strconv.ParseInt(string(b.input[b.offset:endIndex]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer value: %v", err)
	}

	b.offset = endIndex + 1

	return number, nil
}

// parseList parses a bencoded list in the form "l<items>e"
func (b *Bencode) parseList() (any, error) {
	b.offset++ // Skip 'l'
	list := []any{}

	for b.offset < len(b.input) && b.input[b.offset] != 'e' {
		item, err := b.parseNext()
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if b.offset == len(b.input) || b.input[b.offset] != 'e' {
		return nil, errors.New("invalid list format, missing 'e'")
	}

	b.offset++ // Skip 'e'
	return list, nil
}

// parseDictionary parses a bencoded dictionary in the form "d<key><value>e"
func (b *Bencode) parseDictionary() (any, error) {
	b.offset++ // Skip 'd'
	dict := map[string]any{}
	for b.offset < len(b.input) && b.input[b.offset] != 'e' {
		// Parse the key, which is always a string
		key, err := b.parseString()
		if err != nil {
			return nil, err
		}
		// Parse the value associated with the key
		value, err := b.parseNext()
		if err != nil {
			return nil, err
		}
		dict[string(key.([]byte))] = value
	}

	if b.offset == len(b.input) || b.input[b.offset] != 'e' {
		return nil, errors.New("invalid dictionary format, missing 'e'")
	}

	b.offset++ // Skip 'e'
	return dict, nil
}

// parseNext identifies the next element and delegates to the correct parser.
func (b *Bencode) parseNext() (any, error) {
	if b.offset >= len(b.input) {
		return nil, errors.New("unexpected end of input")
	}

	switch b.input[b.offset] {
	case 'i':
		return b.parseInt()
	case 'l':
		return b.parseList()
	case 'd':
		return b.parseDictionary()
	default:
		if b.input[b.offset] >= '0' && b.input[b.offset] <= '9' {
			return b.parseString()
		}
		return nil, fmt.Errorf("unknown token: %c", b.input[b.offset])
	}
}

// Parse initializes the parsing process.
func (b *Bencode) Parse() error {
	for b.offset < len(b.input) {
		item, err := b.parseNext()
		if err != nil {
			return err
		}
		b.Result = append(b.Result, item)
	}

	return nil
}

// New creates a new Bencode instance.
func New(input []byte) Bencode {
	return Bencode{
		input: input,
	}
}
