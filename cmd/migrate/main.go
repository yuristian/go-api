package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/yuristian/go-api/internal/config"
)

func buildDBURL(cfg *config.Config) (string, error) {
	switch strings.ToLower(cfg.DB.Type) {
	case "postgres", "postgresql":
		// postgres://user:pass@host:port/dbname?sslmode=disable
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode), nil
	case "mysql":
		// mysql://user:pass@tcp(host:port)/dbname
		return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s",
			cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name), nil
	default:
		return "", fmt.Errorf("unsupported driver: %s", cfg.DB.Type)
	}
}

func createMigrationFiles(name string) error {
	if name == "" {
		return fmt.Errorf("migration name is required")
	}
	ts := time.Now().Unix()
	base := fmt.Sprintf("%d_%s", ts, strings.ReplaceAll(name, " ", "_"))
	up := path.Join("migrations", base+".up.sql")
	down := path.Join("migrations", base+".down.sql")

	if err := os.WriteFile(up, []byte("-- up migration\n\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(down, []byte("-- down migration\n\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("Created:\n - %s\n - %s\n", up, down)
	return nil
}

func main() {
	action := flag.String("action", "up", "migration action: create|up|down|version|force")
	name := flag.String("name", "", "name for create (used with -action=create)")
	steps := flag.Int("steps", 0, "number of steps for down (optional)")
	force := flag.Int("version", 0, "force set version (used with -action=force)")
	flag.Parse()

	// load config
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("failed to load config")
	}

	// support create locally without migrate lib
	if *action == "create" {
		if err := createMigrationFiles(*name); err != nil {
			log.Fatalf("create failed: %v", err)
		}
		return
	}

	dbURL, err := buildDBURL(cfg)
	if err != nil {
		log.Fatalf("build db url: %v", err)
	}

	m, err := migrate.New("file://./migrations", dbURL)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("m.Up: %v", err)
		}
		fmt.Println("migrations up applied")
	case "down":
		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("m.Steps(-n): %v", err)
			}
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("m.Down: %v", err)
			}
		}
		fmt.Println("migrations down executed")
	case "version":
		v, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("m.Version: %v", err)
		}
		fmt.Printf("current version: %d, dirty: %v\n", v, dirty)
	case "force":
		if err := m.Force(*force); err != nil {
			log.Fatalf("m.Force: %v", err)
		}
		fmt.Printf("forced to version %d\n", *force)
	default:
		log.Fatalf("unknown action: %s", *action)
	}
}
