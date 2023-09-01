CREATE TABLE IF NOT EXISTS sessions_ping (
	session_id	UUID REFERENCES sessions(id),
	date_time	INT,
    createdAt   TIMESTAMP DEFAULT current_timestamp
);