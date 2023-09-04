CREATE TABLE IF NOT EXISTS activity (
	session_id		UUID REFERENCES sessions(id),
	session_type	VARCHAR(20),
	login			VARCHAR(50) REFERENCES users(login),
	start_date_time	TIMESTAMP DEFAULT current_timestamp,
	end_date_time	TIMESTAMP DEFAULT current_timestamp,
	CONSTRAINT unique_session_id_type UNIQUE (session_id, session_type)
);