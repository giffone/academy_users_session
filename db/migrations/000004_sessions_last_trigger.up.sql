CREATE OR REPLACE FUNCTION delete_old_sessions()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM sessions_last
    WHERE date_time < (NOW() - INTERVAL '12 hours');
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER sessions_last_delete_trigger
BEFORE SELECT ON sessions_last
FOR EACH STATEMENT
EXECUTE FUNCTION delete_old_sessions();