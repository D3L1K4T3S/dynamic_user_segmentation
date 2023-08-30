package migrate

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"time"
)

const (
	sslDisable    = "?sslmode=disable"
	prefixMigrate = "file://migrations"

	defaultAttempts = 10
	dafaultTimeout  = 1 * time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("Migration: not declared PG_URL")
	}

	databaseURL += sslDisable
	attempts := defaultAttempts
	var mig *migrate.Migrate
	var err error

	for attempts > 0 {
		mig, err = migrate.New(prefixMigrate, databaseURL)
		if err == nil {
			break
		}
		log.Printf("Migrate: Posgtres tryint to connect, attempts left: %d", attempts)
		time.Sleep(dafaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Nigrate: postgres can't connect")
	}

	err = mig.Up()
	defer func() {
		_, _ = mig.Close()
	}()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migrate: %x", migrate.ErrNoChange.Error())
			return
		}
		log.Fatalf("Migrate: %x", migrate.ErrNoChange.Error())
	}

	log.Printf("Migrate: success ended...")
}
