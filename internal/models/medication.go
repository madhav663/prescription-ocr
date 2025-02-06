package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Medication struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Alternatives string                 `json:"alternatives"`
	SideEffects  string                 `json:"side_effects"`
	Analysis     map[string]interface{} `json:"analysis,omitempty"`
	Confidence   float64                 `json:"confidence,omitempty"`
	CreatedAt    *time.Time              `json:"created_at,omitempty"`
	UpdatedAt    *time.Time              `json:"updated_at,omitempty"`
}

type MedicationModel struct {
	DB *sql.DB
}

func (m *MedicationModel) GetMedication(id int) (*Medication, error) {
	if m.DB == nil {
		log.Println("❌ Database connection is nil")
		return nil, fmt.Errorf("database connection is nil")
	}

	log.Printf("🔍 Querying medication with ID: %d", id)

	var medication Medication
	err := m.DB.QueryRow(
		"SELECT id, name, COALESCE(alternatives, ''), COALESCE(side_effects, ''), created_at, updated_at FROM medications WHERE id = $1", id).
		Scan(&medication.ID, &medication.Name, &medication.Alternatives, &medication.SideEffects, &medication.CreatedAt, &medication.UpdatedAt)

	if err != nil {
		log.Printf("❌ Query failed: %v", err)
		return nil, err
	}

	log.Printf("✅ Medication found: %+v", medication)
	return &medication, nil
}

func (m *MedicationModel) AddMedication(medication *Medication) error {
	query := "INSERT INTO medications (name, alternatives, side_effects, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at"
	err := m.DB.QueryRow(query, medication.Name, medication.Alternatives, medication.SideEffects).
		Scan(&medication.ID, &medication.CreatedAt, &medication.UpdatedAt)
	return err
}

func (m *MedicationModel) UpdateMedication(medication *Medication) error {
	query := "UPDATE medications SET name = $1, alternatives = $2, side_effects = $3, updated_at = NOW() WHERE id = $4"
	_, err := m.DB.Exec(query, medication.Name, medication.Alternatives, medication.SideEffects, medication.ID)
	return err
}

func (m *MedicationModel) DeleteMedication(id int) error {
	query := "DELETE FROM medications WHERE id = $1"
	_, err := m.DB.Exec(query, id)
	return err
}
