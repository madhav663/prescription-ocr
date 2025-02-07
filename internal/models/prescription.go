package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

var DB *sql.DB


type Prescription struct {
	ID            int       `json:"id"`
	OriginalImage string    `json:"original_image"`
	ExtractedText string    `json:"extracted_text"`
	CreatedAt     time.Time `json:"created_at"`
}


func SavePrescription(p *Prescription) error {
	query := "INSERT INTO prescriptions (original_image, extracted_text, created_at) VALUES ($1, $2, NOW()) RETURNING id"
	err := DB.QueryRow(query, p.OriginalImage, p.ExtractedText).Scan(&p.ID)
	if err != nil {
		log.Printf("Database Insert Error: %v", err)
		return err
	}
	log.Printf("Prescription stored in DB with ID: %d", p.ID)
	return nil
}


func GetPrescriptions() ([]byte, error) {
	query := "SELECT id, original_image, extracted_text, created_at FROM prescriptions"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []Prescription
	for rows.Next() {
		var p Prescription
		err := rows.Scan(&p.ID, &p.OriginalImage, &p.ExtractedText, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, p)
	}

	jsonData, err := json.Marshal(prescriptions)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
