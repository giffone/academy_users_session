CREATE OR REPLACE FUNCTION session.f_before_insert_check_time_in_campus() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date_time <= NEW.start_date_time THEN
        RAISE EXCEPTION 'end_date_time must be greater than start_date_time';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_before_insert_check_time_in_campus
BEFORE INSERT ON session.in_campus
FOR EACH ROW
EXECUTE FUNCTION session.f_before_insert_check_time_in_campus();


CREATE OR REPLACE FUNCTION session.f_before_update_check_time_in_campus() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date_time <= OLD.end_date_time THEN
        RAISE EXCEPTION 'end_date_time must be greater than previous value';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_before_update_check_time_in_campus
BEFORE UPDATE ON session.in_campus
FOR EACH ROW
EXECUTE FUNCTION session.f_before_update_check_time_in_campus();

GRANT EXECUTE ON FUNCTION session.f_before_insert_check_time_in_campus() TO PUBLIC;
GRANT EXECUTE ON FUNCTION session.f_before_insert_check_time_in_campus() TO session_manager;
GRANT EXECUTE ON FUNCTION session.f_before_insert_check_time_in_campus() TO postgres;
GRANT EXECUTE ON FUNCTION session.f_before_update_check_time_in_campus() TO PUBLIC;
GRANT EXECUTE ON FUNCTION session.f_before_update_check_time_in_campus() TO session_manager;
GRANT EXECUTE ON FUNCTION session.f_before_update_check_time_in_campus() TO postgres;