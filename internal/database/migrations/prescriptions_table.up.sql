CREATE TABLE IF NOT EXISTS prescriptions (
    id SERIAL PRIMARY KEY,
    original_image TEXT NOT NULL,
    extracted_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
