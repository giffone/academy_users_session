CREATE TABLE IF NOT EXISTS sessions (
	session_id		UUID PRIMARY KEY,
	comp_name		VARCHAR(30) REFERENCES computers(comp_name),
	ip_addr			VARCHAR(20),
	login			VARCHAR(50) REFERENCES users(login),
	next_ping_sec	INT,
	start_date_time	TIMESTAMP DEFAULT current_timestamp,
	end_date_time	TIMESTAMP DEFAULT current_timestamp
);