package main

import (
	"bytes"
	"testing"
)

func TestWriteData(t *testing.T) {
	//buffer is a writer
	var b bytes.Buffer
	if err := WriteData(&b); err != nil {
		t.Fail()
	}
	written := b.Bytes()
	t.Log(string(written)) //casting the bytes to a string
}
