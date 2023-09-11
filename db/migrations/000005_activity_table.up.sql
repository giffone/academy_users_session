CREATE TABLE IF NOT EXISTS session.activity (
	session_id		UUID REFERENCES session.in_campus(id),
	session_type	VARCHAR(20),
	login			VARCHAR(50) REFERENCES public.users(login),
	start_date_time	TIMESTAMP DEFAULT current_timestamp,
	end_date_time	TIMESTAMP DEFAULT current_timestamp,
	CONSTRAINT unique_session_id_type UNIQUE (session_id, session_type)
);

ALTER TABLE IF EXISTS session.activity
    OWNER to postgres;

GRANT ALL ON TABLE session.activity TO session_manager;

GRANT ALL ON TABLE session.activity TO postgres;