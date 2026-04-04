CREATE TABLE IF NOT EXISTS new_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    date TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    address TEXT NOT NULL,
    n_of_slots INTEGER NOT NULL,
    n_of_volunteers FLOAT NOT NULL,
    total_hours FLOAT NOT NULL,
    tags TEXT
);

INSERT INTO new_events SELECT * FROM events;
DROP TABLE events;
ALTER TABLE new_events RENAME TO events;