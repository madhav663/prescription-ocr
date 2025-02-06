package schema

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(cfg DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("❌ Failed to open database connection: %v", err)
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		log.Printf("❌ Database ping failed: %v", err)
		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return db, nil
}
