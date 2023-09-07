DROP TRIGGER IF EXISTS t_before_insert_check_time_activity ON activity;
DROP FUNCTION IF EXISTS f_before_insert_check_time_activity();

DROP TRIGGER IF EXISTS t_before_update_check_time_activity ON activity;
DROP FUNCTION IF EXISTS f_before_update_check_time_activity();