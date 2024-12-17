CREATE OR REPLACE FUNCTION create_user(p_email TEXT, p_password TEXT)
RETURNS TABLE (id TEXT, created BOOLEAN) AS $$
DECLARE
    v_id TEXT;
BEGIN
    INSERT INTO users (email, password)
    VALUES (p_email, p_password)
    RETURNING users.id INTO v_id;

    RETURN QUERY SELECT v_id, TRUE;
END;
$$ LANGUAGE plpgsql;