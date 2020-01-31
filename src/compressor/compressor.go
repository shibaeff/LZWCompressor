package compressor

import (
	"fmt"
	"strings"
)

// Compressor provides methods to compresses peices of text into one chunk
type Compressor struct {
	Dictionary map[string]int // Stores the LZW dictionary
	DictSize   int
	Yield      []int // array of compressed values
}

// Compress the given string with LZW
func (comp *Compressor) Compress(uncompressed string) {
	comp.DictSize = 256
	comp.Dictionary = make(map[string]int, comp.DictSize)
	for i := 0; i < comp.DictSize; i++ {
		// encoding to utf8
		comp.Dictionary[string([]byte{byte(i)})] = i
	}

	var result []int
	var w []byte
	for i := 0; i < len(uncompressed); i++ {
		ch := uncompressed[i]
		appended := append(w, ch)
		if _, flag := comp.Dictionary[string(appended)]; flag {
			w = appended
		} else {
			result = append(result, comp.Dictionary[string(w)])
			// Add appended to the dictionary.
			comp.Dictionary[string(appended)] = comp.DictSize
			comp.DictSize++
			//w = []byte{ch}, but re-using appended
			appended[0] = ch
			w = appended[:1]
		}
	}
	if len(w) > 0 {
		result = append(result, comp.Dictionary[string(w)])
	}
	comp.Yield = append(comp.Yield, result...)
}

// Concat concatenates the uncompressed string to the compressed piece of text
func (comp *Compressor) Concat(uncompressed string) {
	dictSize := 256
	var w []byte
	for i := 0; i < len(uncompressed); i++ {
		ch := uncompressed[i]
		appended := append(w, ch)
		if _, ok := comp.Dictionary[string(appended)]; ok {
			w = appended
		} else {
			comp.Yield = append(comp.Yield, comp.Dictionary[string(w)])
			// Add appended to the dictionary.
			comp.Dictionary[string(appended)] = dictSize
			dictSize++
			//w = []byte{ch}, but re-using appended
			appended[0] = ch
			w = appended[:1]
		}
	}

	if len(w) > 0 {
		// Output the code for w.
		comp.Yield = append(comp.Yield, comp.Dictionary[string(w)])
	}
}

// Decompress provides LZW decompression
func (comp *Compressor) Decompress() (string, error) {
	dictSize := 256
	dictionary := make(map[int][]byte, dictSize)
	for i := 0; i < dictSize; i++ {
		dictionary[i] = []byte{byte(i)}
	}

	var result strings.Builder
	var word []byte
	for _, k := range comp.Yield {
		var entry []byte
		if x, ok := dictionary[k]; ok {
			entry = x[:len(x):len(x)]
		} else if k == dictSize && len(word) > 0 {
			entry = append(word, word[0])
		} else {
			return result.String(), UnknownSymbolError(k)
		}
		result.Write(entry)

		if len(word) > 0 {
			word = append(word, entry[0])
			dictionary[dictSize] = word
			dictSize++
		}
		word = entry
	}
	return result.String(), nil
}

//UnknownSymbolError Type serves to mark unknown symbols encountered while decompression
type UnknownSymbolError int

func (e UnknownSymbolError) Error() string {
	return fmt.Sprint("Bad compressed symbol ", int(e))
}
