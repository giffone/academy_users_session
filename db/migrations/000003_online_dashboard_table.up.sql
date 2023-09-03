CREATE TABLE IF NOT EXISTS online_dashboard (
	session_id	UUID REFERENCES sessions(id),
	comp_name	VARCHAR(20) UNIQUE,
	ip_addr		VARCHAR(20),
	login		VARCHAR(50),
	date_time	TIMESTAMP DEFAULT current_timestamp
);