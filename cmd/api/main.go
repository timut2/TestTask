package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/timut2/music-library-api/config"
	_ "github.com/timut2/music-library-api/docs"
	"github.com/timut2/music-library-api/internal/delivery/http"
	"github.com/timut2/music-library-api/internal/repository/api"
	"github.com/timut2/music-library-api/internal/repository/postgresql"
	"github.com/timut2/music-library-api/internal/service"
	"github.com/timut2/music-library-api/pkg/logger"
)

// @title Music Library API
// @version 1.0
// @description This is a sample music library API.

// @host localhost:8080
// @BasePath /
func main() {
	logger.New(os.Stdout, logger.DebugLevel)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := OpenDB(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	songRepo := postgresql.NewSongsRepository(db)
	verseRepo := postgresql.NewMusicInfoRepository(db)
	apiCliet := api.NewApiClient(cfg)

	musicLibrary := service.NewMusicLibrary(songRepo, verseRepo, apiCliet)
	handler := http.NewHandler(musicLibrary)

	srv := NewServer(
		handler,
		cfg,
	)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}

}
