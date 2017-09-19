package agent

const (
	HeadSize = 16
)

type PacHead struct{
	Tag uint32
	Proto uint32
	PayloadLength uint32
}
