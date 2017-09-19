package agent

import (
	"testing"
)

func TestDecBuff(t *testing.T) {
	p := []byte("Hello World")
	EncBuff(p, DefaultEncKey)
	t.Log(string(p))
	DecBuff(p, DefaultEncKey)
	t.Log(string(p))

	if string(p) != "Hello World" {
		t.Error("fail")
	}
}

func TestDecBytes(t *testing.T) {
	p := []byte("Hel")
	EncBytes(p, DefaultEncKey)
	t.Log(string(p))
	DecBytes(p, DefaultEncKey)
	t.Log(string(p))

	if string(p) != "Hel" {
		t.Error("fail")
	}
}
