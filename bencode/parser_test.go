package bencode

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	input := []byte("4:spam")
	bencode := New(input)
	result, err := bencode.parseString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []byte("spam")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseInt(t *testing.T) {
	input := []byte("i42e")
	bencode := New(input)
	result, err := bencode.parseInt()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := 42
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestParseList(t *testing.T) {
	input := []byte("l4:spami42ee")
	bencode := New(input)
	result, err := bencode.parseList()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []interface{}{[]byte("spam"), 42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseDictionary(t *testing.T) {
	input := []byte("d3:cow3:mooe")
	bencode := New(input)
	result, err := bencode.parseDictionary()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := map[string]interface{}{
		"cow": []byte("moo"),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseMultiple(t *testing.T) {
	input := []byte("4:spami42e")
	bencode := New(input)
	err := bencode.Parse()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []interface{}{[]byte("spam"), 42}
	if reflect.DeepEqual(bencode.Result, expected) {
		t.Errorf("Expected %v, got %v", expected, bencode.Result)
	}
}

func TestInvalidString(t *testing.T) {
	input := []byte("4spam")
	bencode := New(input)
	_, err := bencode.parseString()
	if err == nil {
		t.Error("Expected error for invalid string format, got nil")
	}
}

func TestInvalidInt(t *testing.T) {
	input := []byte("i42")
	bencode := New(input)
	_, err := bencode.parseInt()
	if err == nil {
		t.Error("Expected error for missing 'e' in integer, got nil")
	}
}

func TestInvalidList(t *testing.T) {
	input := []byte("l4:spam42e")
	bencode := New(input)
	_, err := bencode.parseList()
	if err == nil {
		t.Error("Expected error for invalid list format, got nil")
	}
}

func TestInvalidDictionary(t *testing.T) {
	input := []byte("d3:cow3:moo")
	bencode := New(input)
	_, err := bencode.parseDictionary()
	if err == nil {
		t.Error("Expected error for invalid dictionary format, got nil")
	}
}
