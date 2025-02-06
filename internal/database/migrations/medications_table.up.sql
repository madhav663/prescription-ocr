
CREATE TABLE IF NOT EXISTS medications (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alternatives TEXT,
    side_effects TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_medications_name ON medications(name);


CREATE OR REPLACE FUNCTION update_medication_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TRIGGER update_medications_timestamp
    BEFORE UPDATE ON medications
    FOR EACH ROW
    EXECUTE FUNCTION update_medication_timestamp();




