package main

import (
	"errors"
	"flag"
	"log"
	"path/filepath"
	"tax-helper/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	action := flag.String("action", "up", "migration action: up or down")
	flag.Parse()

	migrationsPath, err := filepath.Abs("internal/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(
		"file://"+filepath.ToSlash(migrationsPath),
		cfg.DBDsn,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown action: %s, expected 'up' or 'down'", *action)
	}
}
