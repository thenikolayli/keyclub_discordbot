ALTER TABLE events ADD COLUMN n_of_slots INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN tags TEXT DEFAULT "";

ALTER TABLE members ADD COLUMN middle_name TEXT DEFAULT "";

CREATE TABLE IF NOT EXISTS events_leaders (
    event_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    PRIMARY KEY (event_id, member_id),
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);