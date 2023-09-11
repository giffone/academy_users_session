CREATE TABLE IF NOT EXISTS session.in_campus (
	id				UUID PRIMARY KEY,
	comp_name		VARCHAR(30) REFERENCES session.computers(comp_name),
	ip_addr			VARCHAR(20),
	login			VARCHAR(50) REFERENCES public.users(login),
	next_ping_sec	INT,
	start_date_time	TIMESTAMP DEFAULT current_timestamp,
	end_date_time	TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE IF EXISTS session.in_campus
    OWNER to postgres;

GRANT ALL ON TABLE session.in_campus TO session_manager;

GRANT ALL ON TABLE session.in_campus TO postgres;