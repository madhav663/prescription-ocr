
CREATE TABLE medications (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alternatives TEXT,
    side_effects TEXT
);


DROP TABLE IF EXISTS medications;