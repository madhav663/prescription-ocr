package models

import (
    "database/sql"
    "fmt"
    "time"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
    var db *sql.DB
    var err error
    
   
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        db, err = sql.Open("postgres", dataSourceName)
        if err != nil {
            fmt.Printf("Attempt %d: Failed to open database: %v\n", i+1, err)
            time.Sleep(time.Second * 2)
            continue
        }

        err = db.Ping()
        if err != nil {
            fmt.Printf("Attempt %d: Failed to ping database: %v\n", i+1, err)
            time.Sleep(time.Second * 2)
            continue
        }

        
        break
    }

    if err != nil {
        return nil, fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
    }

    
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}