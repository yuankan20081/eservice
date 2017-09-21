package license

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestManager_Watch(t *testing.T) {
	m := NewManager()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	if err := m.Watch(ctx, "./license"); err != nil {
		t.Error(err)
	}
}
