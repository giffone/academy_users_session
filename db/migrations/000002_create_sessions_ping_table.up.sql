CREATE TABLE IF NOT EXISTS sessions_ping (
	session_id		UUID REFERENCES sessions(id),
	session_type	VARCHAR(15),
	next_ping_date	TIMESTAMP DEFAULT current_timestamp,
    created			TIMESTAMP DEFAULT current_timestamp,
	updated			TIMESTAMP,
	CONSTRAINT unique_session_id_type UNIQUE (session_id, session_type)
);