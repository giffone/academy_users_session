DROP TRIGGER IF EXISTS sessions_last_delete_trigger ON sessions_last;

DROP FUNCTION IF EXISTS delete_old_sessions();