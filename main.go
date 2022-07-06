package main

import (
	"embed"
	"io/fs"
	"runtime"

	"github.com/BililiveRecorder/BililiveRecorder-Tray/cmd/tray"
)

//go:embed frontend/dist/*
//go:embed frontend/dist/assets/*
var assets embed.FS

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	assets, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		panic(err)
	}
	tray.Main(assets)
}
