
DROP TRIGGER IF EXISTS update_medications_timestamp ON medications;


DROP FUNCTION IF EXISTS update_medication_timestamp();


DROP INDEX IF EXISTS idx_medications_name;

DROP TABLE IF EXISTS medications;