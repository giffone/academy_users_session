CREATE TABLE IF NOT EXISTS sessions_ping (
	session_id		UUID REFERENCES sessions(id),
	session_type	VARCHAR(20),
	date_time		TIMESTAMP DEFAULT current_timestamp,
	CONSTRAINT unique_session_id_type UNIQUE (session_id, session_type)
);