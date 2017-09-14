package game_net

const (
	MaxPayloadLength = 64 * 1024
)

type PacketHead struct {
	Tag           uint32
	Proto         uint32
	PayloadLength uint32
}
