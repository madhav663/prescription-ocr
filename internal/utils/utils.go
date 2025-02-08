// utils.go

package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/madhav663/prescription-ocr/internal/database/schema"
)

// connectWithRetry tries to connect to the database with retries.
func ConnectWithRetry(dbConfig schema.DBConfig, retries int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		log.Printf("🛠️ Attempting to connect to DB (Attempt %d/%d)...", i+1, retries)
		db, err = schema.NewDatabase(dbConfig)
		if err == nil {
			log.Println("✅ Database connection successful.")
			return db, nil
		}
		log.Printf("⚠️ Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
	}
	return nil, fmt.Errorf("database not reachable after %d attempts: %v", retries, err)
}
