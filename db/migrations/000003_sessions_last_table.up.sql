CREATE TABLE IF NOT EXISTS sessions_last (
	id			UUID,
	comp_name	VARCHAR(15) UNIQUE,
	ip_addr		VARCHAR(15),
	login		VARCHAR(50),
	date_time	TIMESTAMP DEFAULT current_timestamp
);