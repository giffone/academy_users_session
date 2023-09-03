CREATE TABLE IF NOT EXISTS online_dashboard (
	session_id		UUID REFERENCES sessions(id),
	comp_name		VARCHAR(15) UNIQUE,
	ip_addr			VARCHAR(15),
	login			VARCHAR(50),
	next_ping_date	TIMESTAMP DEFAULT current_timestamp,
	created		   	TIMESTAMP DEFAULT current_timestamp,
	updated   		TIMESTAMP
);