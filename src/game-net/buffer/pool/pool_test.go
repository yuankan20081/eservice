package pool

import (
	"io"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	p := New()

	buf := p.Get()
	buf.Reset()
	buf.SetProto(111)
	buf.Write([]byte("hello buffer"))
	if n, err := io.Copy(os.Stdout, buf); err != nil {
		t.Error(err)
	} else {
		t.Log(n)
	}
	p.Put(buf)
}
