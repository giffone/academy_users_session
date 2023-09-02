CREATE OR REPLACE FUNCTION delete_old_sessions()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM sessions_last
    WHERE date_time < EXTRACT(EPOCH FROM NOW() - INTERVAL '10 minutes');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER sessions_last_delete_trigger
BEFORE INSERT OR UPDATE ON sessions_last
FOR EACH ROW
EXECUTE FUNCTION delete_old_sessions();
