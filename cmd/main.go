package main

import (
	"context"
	"musictags-joiner/internal/artists"
	"musictags-joiner/internal/genres"
	"musictags-joiner/pkgs/logger"
	"musictags-joiner/pkgs/storage"
	"os"
	"sync"
)

const (
	serviceName string = "musictags-joiner"
)

var (
	version    = "No Version Provided"
	commitHash = "No Git Commit Hash Provided"
)

type application struct {
	logger *logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
	cfg    config
	//services
	srvGenres  *genres.Service
	srvArtists *artists.Service
	// storage
	sync.RWMutex
	storage *storage.Storage
}

const (
	stortagsPath string = "./stortags/"
	logFilePath  string = "./logs"
)

func main() {
	app := application{}
	if err := app.initConfig(""); err != nil {
		os.Exit(1)
	}
	app.logger = logger.Init(app.cfg.LogLevel, app.cfg.LogFile)

	app.storage = storage.NewStorage(app.cfg.Apikey, stortagsPath)
	app.srvGenres = genres.NewService(app.storage, app.logger, stortagsPath)
	app.srvArtists = artists.NewService(app.storage, app.logger, stortagsPath)

	if _, err := os.Stat(stortagsPath); os.IsNotExist(err) {
		os.Mkdir(stortagsPath, 0755)
	}

	app.ctx, app.cancel = context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go app.storage.StorageBouncer(&wg, app.ctx)

	app.runServer()
	wg.Wait()
}
