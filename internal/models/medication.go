package models

import (
	"database/sql"
)


type Medication struct {
	ID           int
	Name         string
	Alternatives string
	SideEffects  string
}

type MedicationModel struct {
	DB *sql.DB
}


func (m *MedicationModel) GetMedication(id int) (*Medication, error) {
	query := "SELECT id, name, alternatives, side_effects FROM medications WHERE id = $1"
	row := m.DB.QueryRow(query, id)

	medication := &Medication{}
	err := row.Scan(&medication.ID, &medication.Name, &medication.Alternatives, &medication.SideEffects)
	if err != nil {
		return nil, err
	}

	return medication, nil
}

func (m *MedicationModel) AddMedication(medication *Medication) error {
	query := "INSERT INTO medications (name, alternatives, side_effects) VALUES ($1, $2, $3) RETURNING id"
	err := m.DB.QueryRow(query, medication.Name, medication.Alternatives, medication.SideEffects).Scan(&medication.ID)
	return err
}

func (m *MedicationModel) UpdateMedication(medication *Medication) error {
	query := "UPDATE medications SET name = $1, alternatives = $2, side_effects = $3 WHERE id = $4"
	_, err := m.DB.Exec(query, medication.Name, medication.Alternatives, medication.SideEffects, medication.ID)
	return err
}

func (m *MedicationModel) DeleteMedication(id int) error {
	query := "DELETE FROM medications WHERE id = $1"
	_, err := m.DB.Exec(query, id)
	return err
}