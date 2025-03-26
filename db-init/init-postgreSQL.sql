CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    barcode CHAR(13) UNIQUE NOT NULL,
    contact VARCHAR(255),
    remaining_entries INT DEFAULT 0 CHECK (remaining_entries >= 0),
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_entries INT DEFAULT 0 CHECK (total_entries >= 0)
);
