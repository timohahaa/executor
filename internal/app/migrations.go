package app

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	attempts = 20
	timeout  = time.Second
)

func init() {
	dbURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(dbURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}

	dbURL += "?sslmode=disable"

	var m *migrate.Migrate
	var err error

	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbURL)
		if err == nil {
			break
		}
		log.Printf("migrations: trying to connect, attempts left: %d", attempts)
		time.Sleep(timeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("migrations: connect error: %s", err)
	}

	err = m.Up()
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrations: migrating up error: %s", err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrations: no change")
		return
	}

	log.Printf("migrations: migrationg up success")
}
