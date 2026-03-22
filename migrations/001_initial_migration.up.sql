CREATE TABLE IF NOT EXISTS members (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	nickname TEXT,
	term_hours FLOAT NOT NULL,
	all_hours FLOAT NOT NULL,
	shirt_size TEXT,
	paid_dues BOOLEAN,
	grad_year INTEGER,
	strikes INTEGER
)