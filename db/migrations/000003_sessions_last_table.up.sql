CREATE TABLE IF NOT EXISTS sessions_last (
	session_id	UUID REFERENCES sessions(id),
	comp_name	VARCHAR(15),
	ip_addr		VARCHAR(15),
	login		VARCHAR(50),
	date_time	TIMESTAMP DEFAULT current_timestamp,
	createdAt	TIMESTAMP DEFAULT current_timestamp
);