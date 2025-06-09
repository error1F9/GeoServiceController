package main

import (
	"GeoService/config"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"strconv"
)

var (
	command    = flag.String("cmd", "up", "goose command (up, down, status, etc.)")
	migrations = flag.String("dir", "migrations", "directory with migration files")
)

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.NewAppConfig()

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	db, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatal(err)
	}

	switch *command {
	case "up":
		if err := goose.Up(db, *migrations); err != nil {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
	case "down":
		if err := goose.Down(db, *migrations); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
	case "down-to":
		if flag.NArg() < 1 {
			log.Fatal("Missing version argument for down-to command")
		}
		version, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		if err := goose.DownTo(db, *migrations, int64(version)); err != nil {
			log.Fatalf("Failed to rollback to version %s: %v", version, err)
		}
	case "status":
		if err := goose.Status(db, *migrations); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	case "version":
		if err := goose.Version(db, *migrations); err != nil {
			log.Fatalf("Failed to get current version: %v", err)
		}
	default:
		log.Printf("Unknown command: %q", *command)
		os.Exit(1)
	}

}
