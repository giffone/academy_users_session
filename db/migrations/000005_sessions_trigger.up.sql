CREATE OR REPLACE FUNCTION f_before_insert_check_time_sessions() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date_time <= NEW.start_date_time THEN
        RAISE EXCEPTION 'end_date_time must be greater than start_date_time';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_before_insert_check_time_sessions
BEFORE INSERT ON sessions
FOR EACH ROW
EXECUTE FUNCTION f_before_insert_check_time_sessions();


CREATE OR REPLACE FUNCTION f_before_update_check_time_sessions() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date_time <= OLD.end_date_time THEN
        RAISE EXCEPTION 'end_date_time must be greater than previous value';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_before_update_check_time_sessions
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION f_before_update_check_time_sessions();