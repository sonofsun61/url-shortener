package main

import (
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/lib/logger"
	"urlShortener/internal/lib/logger/sl"
	"urlShortener/internal/storage/sqlite"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	var config config.Config
	config.MustNewConfig()

	logger := logger.SetupLogger(config.Env)
	logger = logger.With(slog.String("env", config.Env))
	logger.Info("Initializing server", slog.String("address", config.Address))
	logger.Debug("Logger debug module 	enabled")

	var storage sqlite.Storage = sqlite.Storage{}
	err := storage.NewStorage(config.StoragePath)
	if err != nil {
		logger.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// id, err := storage.SaveURL("https://www.apple.com", "apple")
	// if err != nil {
	// 	logger.Error("Failed to save URL", sl.Err(err))
	// 	os.Exit(1)
	// }
	// logger.Info("Saved URL", slog.Int64("id", id))
	// fmt.Println(id)

	// url, err := storage.GetURL("tbank")
	// if err != nil {
	// 	logger.Error("Failed to get URL", sl.Err(err))
	// }
	// fmt.Println(url)

	err = storage.DeleteURL("tbank")
	if err != nil {
		logger.Error("Failed to delete URL", sl.Err(err))
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	//TODO: run server
}
