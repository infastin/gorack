package validation_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/infastin/gorack/validation"
	isint "github.com/infastin/gorack/validation/is/int"
	isstr "github.com/infastin/gorack/validation/is/str"
)

type Config struct {
	Logger   LoggerConfig
	Database DatabaseConfig
}

func (cfg *Config) Validate() error {
	return validation.All(
		validation.Ptr(&cfg.Logger, "logger").With(validation.Custom),
		validation.Ptr(&cfg.Database, "database").With(validation.Custom),
	)
}

type LoggerConfig struct {
	Level string
}

func (cfg *LoggerConfig) Validate() error {
	return validation.All(
		validation.String(cfg.Level, "level").Required(true).In("debug", "info", "warn", "error"),
	)
}

type DatabaseConfig struct {
	Backend string

	Postgres PostgresConfig
	SQLite   SQLiteConfig
}

func (cfg *DatabaseConfig) Validate() error {
	return validation.All(
		validation.String(cfg.Backend, "backend").Required(true).In("postgres", "sqlite"),
		validation.PtrI(&cfg.Postgres).If(cfg.Backend == "postgres").With(validation.Custom).EndIf(),
		validation.PtrI(&cfg.SQLite).If(cfg.Backend == "sqlite").With(validation.Custom).EndIf(),
	)
}

type PostgresConfig struct {
	Host string
	Port int
}

func (cfg *PostgresConfig) Validate() error {
	return validation.All(
		validation.String(cfg.Host, "host").Required(true).With(isstr.Host),
		validation.Number(cfg.Port, "port").Required(true).With(isint.Port),
	)
}

type SQLiteConfig struct {
	Path string
}

func (cfg *SQLiteConfig) Validate() error {
	return validation.All(
		validation.String(cfg.Path, "path").Required(true),
	)
}

func Example_config() {
	sqlite := DatabaseConfig{
		Backend: "sqlite",
	}
	fmt.Println(sqlite.Validate())

	postgres := DatabaseConfig{
		Backend: "postgres",
	}
	fmt.Println(postgres.Validate())

	config := Config{
		Logger: LoggerConfig{
			Level: "panic",
		},
		Database: DatabaseConfig{
			Backend: "postgres",
		},
	}

	b, _ := json.MarshalIndent(config.Validate(), "", "  ")
	os.Stdout.Write(b)
}
