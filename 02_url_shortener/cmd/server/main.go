package main

import (
	"crypto/tls"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"url-shortener/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger *slog.Logger
	urls   models.URLModelInterface
}

func main() {

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))

	dir, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Debug("initializing app", "working dir", dir)

	db, err := openDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		urls:   &models.URLModel{DB: db},
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		},
	}

	srv := http.Server{
		Addr:         ":8989",
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", slog.String("addr", srv.Addr))

	if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:./data/hurl.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
