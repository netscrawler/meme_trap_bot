package storage

const imagesTable = `
CREATE TABLE IF NOT EXISTS images (
	id INTEGER PRIMARY KEY AUTOINCREMENT,

    file_id TEXT NOT NULL,
    file_unique_id TEXT NOT NULL UNIQUE,
    file_size INTEGER,
    width INTEGER,
    height INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by TEXT NOT NULL,
    added_by_id INTEGER NOT NULL
);
`
