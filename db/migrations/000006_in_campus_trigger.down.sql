DROP TRIGGER IF EXISTS t_before_insert_check_time_in_campus ON session.in_campus;
DROP FUNCTION IF EXISTS session.f_before_insert_check_time_in_campus();

DROP TRIGGER IF EXISTS t_before_update_check_time_in_campus ON session.in_campus;
DROP FUNCTION IF EXISTS session.f_before_update_check_time_in_campus();