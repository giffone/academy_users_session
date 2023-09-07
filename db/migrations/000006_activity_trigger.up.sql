CREATE OR REPLACE FUNCTION f_before_insert_check_time_activity() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date_time <= NEW.start_date_time THEN
        RAISE EXCEPTION 'end_date_time must be greater than start_date_time';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_before_insert_check_time_activity
BEFORE INSERT ON activity
FOR EACH ROW
EXECUTE FUNCTION f_before_insert_check_time_activity();


CREATE OR REPLACE FUNCTION f_before_update_check_time_activity() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date_time <= OLD.end_date_time THEN
        RAISE EXCEPTION 'end_date_time must be greater than previous value';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_before_update_check_time_activity
BEFORE UPDATE ON activity
FOR EACH ROW
EXECUTE FUNCTION f_before_update_check_time_activity();