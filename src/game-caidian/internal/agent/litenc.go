package agent

import "encoding/binary"

type TNACC struct {
	B [4]byte
}

func DecBytes(p []byte, k uint32) {
	var acc TNACC
	binary.LittleEndian.PutUint32(acc.B[:], k)
	h := acc.B[0] ^ acc.B[1] ^ acc.B[2] ^ acc.B[3]
	l := len(p)

	for i := 0; i < l; i++ {
		p[i] ^= h
	}
}

func EncBytes(p []byte, k uint32) {
	DecBytes(p, k)
}

func EncBuff(p []byte, k uint32) {
	length := len(p)
	for length >= 4 {
		v := binary.LittleEndian.Uint32(p[:4])
		binary.LittleEndian.PutUint32(p[:4], ((v<<17)|((v>>15)&0x1ffff))^k)
		length -= 4
		p = p[4:]
	}

	if length != 0 {
		EncBytes(p, k)
	}
}

func DecBuff(p []byte, k uint32) {
	length := len(p)
	for length >= 4 {
		v := binary.LittleEndian.Uint32(p[:4]) ^ k
		binary.LittleEndian.PutUint32(p[:4], (v<<15)|((v>>17)&0x7fff))
		length -= 4
		p = p[4:]
	}

	if length != 0 {
		DecBytes(p, k)
	}
}
