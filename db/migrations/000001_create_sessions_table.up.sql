CREATE TABLE IF NOT EXISTS sessions (
	id			UUID PRIMARY KEY,
	comp_name	VARCHAR(15),
	ip_addr		VARCHAR(15),
	login		VARCHAR(50),
	status		VARCHAR(15),
	date_time	INT,
	createdAt	TIMESTAMP DEFAULT current_timestamp
);