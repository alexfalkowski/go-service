package test

import (
	"github.com/golang-migrate/migrate/v4"
	// This are need to use migrate in test.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// NewMigrator for test.
func NewMigrator() *Migrator {
	m, _ := migrate.New("file://../../../test/migrations", NewPGConfig().URL)

	return &Migrator{migrate: m}
}

// Migrator for test.
type Migrator struct {
	migrate *migrate.Migrate
}

// Up to latest version.
func (m *Migrator) Up() error {
	return m.migrate.Up()
}

// Drop the DB.
func (m *Migrator) Drop() error {
	return m.migrate.Drop()
}
