CREATE TABLE IF NOT EXISTS session.computers (
	comp_name   VARCHAR(30) PRIMARY KEY,
    status      VARCHAR(15)
);

ALTER TABLE IF EXISTS session.computers
    OWNER to postgres;

GRANT ALL ON TABLE session.computers TO session_manager;

GRANT ALL ON TABLE session.computers TO postgres;