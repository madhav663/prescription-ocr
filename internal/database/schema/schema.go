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
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		log.Fatal("❌ Database configuration is missing! Check environment variables.")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	log.Printf("🔍 Connecting to PostgreSQL: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("❌ Failed to open database connection: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("❌ Database ping failed: %v", err)
		return nil, err
	}

	log.Println("✅ Database connection established successfully")
	return db, nil
}
