package main

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type application struct {
	logger *logInterface
	sync.RWMutex
	cfg     config
	storage storage
	ctx     context.Context
	cancel  context.CancelFunc
}

type storage struct {
	sync.RWMutex
	cache map[string]*mycache
}

const (
	StortagsPath string        = "./stortags/"
	logFilePath  string        = "./logs"
	UrlApiBase   string        = "http://ws.audioscrobbler.com/2.0/"
	cachettl     time.Duration = 1 * time.Hour

	MethodTagGetTopArtists string = "tag.gettopartists"
	MethodTagGetTopTags    string = "tag.getTopTags"
)

var (
	json = jsoniter.ConfigFastest
	r    = strings.NewReplacer(" ", "_")
)

func main() {
	app := application{}
	if err := app.initConfig(""); err != nil {
		os.Exit(1)
	}
	app.storage.cache = make(map[string]*mycache)

	if _, err := os.Stat(StortagsPath); os.IsNotExist(err) {
		os.Mkdir(StortagsPath, 0755)
	}

	app.ctx, app.cancel = context.WithCancel(context.Background())
	app.logger = initLogger(&app.cfg.LogLevel, &app.cfg.LogFile)

	var wg sync.WaitGroup
	wg.Add(1)
	go app.storageBouncer(&wg)

	app.runServer()

	wg.Wait()
	app.logger.Output.Close()
}
