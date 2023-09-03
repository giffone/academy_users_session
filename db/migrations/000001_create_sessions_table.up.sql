CREATE TABLE IF NOT EXISTS sessions (
	id				UUID PRIMARY KEY,
	comp_name		VARCHAR(20),
	ip_addr			VARCHAR(20),
	login			VARCHAR(50),
	next_ping_sec	INT,
	date_time		TIMESTAMP DEFAULT current_timestamp
);