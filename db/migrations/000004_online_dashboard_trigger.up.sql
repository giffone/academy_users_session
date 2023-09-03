CREATE OR REPLACE FUNCTION delete_old_sessions()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM online_dashboard
    WHERE date_time < (NOW() - INTERVAL '12 hours');
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER online_dashboard_delete_trigger
BEFORE SELECT ON online_dashboard
FOR EACH STATEMENT
EXECUTE FUNCTION delete_old_sessions();