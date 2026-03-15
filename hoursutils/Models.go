package hoursutils

// struct to represent a member
type Member struct {
	ID        int     `db:"id"`
	Name      string  `db:"name"`
	Nickname  string  `db:"nickname"`
	TermHours float64 `db:"term_hours"`
	AllHours  float64 `db:"all_hours"`
	ShirtSize string  `db:"shirt_size"`
	PaidDues  bool    `db:"paid_dues"`
	GradYear  int     `db:"grad_year"`
	Strikes   int     `db:"strikes"`
}

var MemberSchema = `
CREATE TABLE IF NOT EXISTS members (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	nickname TEXT,
	term_hours FLOAT NOT NULL,
	all_hours FLOAT NOT NULL,
	shirt_size TEXT,
	paid_dues BOOLEAN,
	grad_year INTEGER,
	strikes INTEGER
)
`
