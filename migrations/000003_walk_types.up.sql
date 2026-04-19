-- Add walk types and bind slots to walk type
CREATE TABLE IF NOT EXISTS walk_types (
    id SERIAL PRIMARY KEY,
    instructor_id INTEGER NOT NULL REFERENCES instructors(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    max_people INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE slots
    ADD COLUMN IF NOT EXISTS walk_type_id INTEGER REFERENCES walk_types(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_walk_types_instructor_id ON walk_types(instructor_id);
