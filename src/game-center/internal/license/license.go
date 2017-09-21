package license

import (
	"fmt"
	"time"
)

type License struct {
	Dir    string
	Server string
	BindIP string
	Expire time.Time
}

func (l *License) String() string {
	return fmt.Sprintf("<%s> <%s> <%s>", l.Server, l.BindIP, l.Expire)
}
