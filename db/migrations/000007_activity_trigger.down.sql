DROP TRIGGER IF EXISTS t_before_insert_check_time_activity ON session.activity;
DROP FUNCTION IF EXISTS session.f_before_insert_check_time_activity();

DROP TRIGGER IF EXISTS t_before_update_check_time_activity ON session.activity;
DROP FUNCTION IF EXISTS session.f_before_update_check_time_activity();