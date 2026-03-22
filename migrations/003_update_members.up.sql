DROP TABLE members;

CREATE TABLE IF NOT EXISTS members (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	nickname TEXT,
    all_hours FLOAT NOT NULL,
	term_hours FLOAT NOT NULL,
    class_year TEXT,
    grad_year INTEGER,
    strikes INTEGER,
    personal_email TEXT,
    school_email TEXT,
    phone_number TEXT,
	shirt_size TEXT,
	paid_dues BOOLEAN
);