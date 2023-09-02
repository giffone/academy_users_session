CREATE TABLE IF NOT EXISTS sessions_last (
	session_id	UUID REFERENCES sessions(id),
	comp_name	VARCHAR(15),
	ip_addr		VARCHAR(15),
	login		VARCHAR(50),
	date_time	INT,
	createdAt	TIMESTAMP DEFAULT current_timestamp
);

CREATE OR REPLACE FUNCTION delete_old_sessions()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM sessions_last
    WHERE date_time < (EXTRACT(EPOCH FROM current_timestamp) - 600)::INT;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_sessions_last
BEFORE INSERT ON sessions_last
FOR EACH ROW
EXECUTE FUNCTION delete_old_sessions();
