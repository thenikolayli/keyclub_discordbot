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
    leaders TEXT NOT NULL DEFAULT "",
    made_by TEXT NOT NULL DEFAULT "",
    sign_up_url TEXT NOT NULL DEFAULT ""
);

INSERT INTO new_events (
    id,
    name,
    date,
    start_time,
    end_time,
    address,
    n_of_slots,
    n_of_volunteers,
    total_hours,
    leaders,
    made_by,
    sign_up_url
)
SELECT
    id,
    name,
    date,
    start_time,
    end_time,
    address,
    n_of_slots,
    n_of_volunteers,
    total_hours,
    "",
    "",
    ""
FROM events;

DROP TABLE events;
ALTER TABLE new_events RENAME TO events;
