package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/xsadia/secred/pkg/utils"
)

const (
	defaultMaxIdleConns = 5
	defaultMaxOpenConns = 10
	defaultMaxLifetime  = 5
)

type Storage struct {
	db            *sql.DB
	migrationPath string
	init          bool
}

var lock = &sync.Mutex{}
var storage *Storage

func NewDB(host, user, password, name, ssl, migrationPath string) (*Storage, error) {
	connectionString :=
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=%s",
			utils.Or(os.Getenv("DB_HOST"), host),
			utils.Or(os.Getenv("DB_USER"), user),
			utils.Or(os.Getenv("DB_PASSWORD"), password),
			utils.Or(os.Getenv("DB_NAME"), name),
			utils.Or(os.Getenv("DB_SSL_MODE"), ssl),
		)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	storage = &Storage{
		db:            db,
		init:          false,
		migrationPath: utils.Or(migrationPath, "file://./migrations"),
	}

	if err := storage.config(); err != nil {
		storage.db.Close()
		return nil, err
	}

	return storage, nil
}

func Close() {
	if storage == nil {
		return
	}

	storage.db.Close()
}

func GetInstance() (*Storage, error) {
	if storage != nil {
		return storage, nil
	}

	lock.Lock()
	defer lock.Unlock()

	newStorage, err := NewDB("localhost",
		"root",
		"root",
		"secred",
		"disable",
		"",
	)
	if err != nil {
		return nil, err
	}

	storage = newStorage
	return storage, nil
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) config() error {
	if s.init {
		return nil
	}

	s.init = true
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

	s.db.SetMaxIdleConns(maxIdleConns)
	s.db.SetMaxOpenConns(maxOpenConns)
	s.db.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))

	return s.migrate()
}

func (s *Storage) SetMigrationPath(path string) {
	s.migrationPath = path
}

func (s *Storage) migrate() error {
	driver, err := postgres.WithInstance(storage.db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		s.migrationPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
