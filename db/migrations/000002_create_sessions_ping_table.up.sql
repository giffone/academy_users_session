CREATE TABLE IF NOT EXISTS sessions_ping (
	session_id		UUID UNIQUE REFERENCES sessions(id),
	session_type	VARCHAR(15),
	date_time		TIMESTAMP DEFAULT current_timestamp,
    createdAt   	TIMESTAMP DEFAULT current_timestamp,
	updatedAt   	TIMESTAMP
);