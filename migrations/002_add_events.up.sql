ALTER TABLE members ADD COLUMN phone_number TEXT NOT NULL DEFAULT '';
ALTER TABLE members ADD COLUMN personal_email TEXT NOT NULL DEFAULT '';
ALTER TABLE members ADD COLUMN school_email TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    address TEXT NOT NULL,
    n_of_volunteers INTEGER NOT NULL,
    total_hours FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS events_members (
    event_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    PRIMARY KEY (event_id, member_id),
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);