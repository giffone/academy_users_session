DROP TRIGGER IF EXISTS t_before_insert_check_time_sessions ON sessions;
DROP FUNCTION IF EXISTS f_before_insert_check_time_sessions();

DROP TRIGGER IF EXISTS t_before_update_check_time_sessions ON sessions;
DROP FUNCTION IF EXISTS f_before_update_check_time_sessions();