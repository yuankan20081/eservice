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
		msg := fmt.Sprintf(format, args...)
		log.Println("[DEBUGIN]", msg)
	}
}
