package compressor

import (
	"fmt"
	"testing"
)

func TestCompress(t *testing.T) {
	var c Compressor
	c.Compress("hello")
	if fmt.Sprint(c.Yield) != "[104 101 108 108 111]" {
		t.Errorf("Got slice %v instead of %v", c.Yield, "[104 101 108 108 111]")
	}
}

func TestConcat(t *testing.T) {
	var c Compressor
	c.Compress("h")
	c.Concat("ello")
	if fmt.Sprint(c.Yield) != "[104 101 108 108 111]" {
		t.Errorf("Got slice %v instead of %v", c.Yield, "[104 101 108 108 111]")
	}
}

func TestCompressor_Decompress(t *testing.T) {
	var c Compressor
	c.Compress("h")
	c.Concat("ello")
	res, err := c.Decompress()
	if err != nil {
		t.Error("Bad symbol encountered")
	}
	if res != "hello" {
		t.Error("Can't decompress hello")
	}
}
