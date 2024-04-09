package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/xsadia/secred/pkg/utils"
)

func New(host, user, password, name, sll string) (*sql.DB, error) {
	connectionString :=
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=%s",
			utils.Or(os.Getenv("DB_HOST"), host),
			utils.Or(os.Getenv("DB_USER"), user),
			utils.Or(os.Getenv("DB_PASSWORD"), password),
			utils.Or(os.Getenv("DB_NAME"), name),
			utils.Or(os.Getenv("DB_SLL_MODE"), sll),
		)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

const (
	defaultMaxIdleConns = 5
	defaultMaxOpenConns = 10
	defaultMaxLifetime  = 5
)

func ConfigDB(db *sql.DB) {
	maxIdleConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	if err != nil {
		maxIdleConns = defaultMaxIdleConns
	}

	maxOpenConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	if err != nil {
		maxOpenConns = defaultMaxOpenConns
	}

	maxLifetime, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_LIFETIME"))
	if err != nil {
		maxLifetime = defaultMaxLifetime
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))
}

func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
