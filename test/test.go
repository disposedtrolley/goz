package test

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type Gamefile string

const (
	JigsawZ8 Gamefile = "gamefiles/Jigsaw.z8"
	JigsawZ5 Gamefile = "gamefiles/Jigsaw.z5"
	ZorkZ3   Gamefile = "gamefiles/zork1.z5" // Yes the extension is .z5, but Zork is a Z3 game
)

func ReadGamfile(g Gamefile) ([]byte, error) {
	_, fname, _, _ := runtime.Caller(0)
	return ioutil.ReadFile(filepath.Join(filepath.Dir(fname), string(g)))
}
