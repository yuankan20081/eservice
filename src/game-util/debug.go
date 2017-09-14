package game_util

import (
	"fmt"
	"log"
)

type BuildMode int

const (
	DebugMode BuildMode = iota
	ReleaseMode
)

var (
	defaultMode BuildMode = DebugMode
)

func SetMode(mode BuildMode) {
	defaultMode = mode
}

func Debug(format string, args ...interface{}) {
	if defaultMode == DebugMode {
		log.Println(fmt.Sprintf(format, args...))
	}
}
