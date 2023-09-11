CREATE TABLE IF NOT EXISTS public.users (
	login   VARCHAR(50) PRIMARY KEY,
    status  VARCHAR(15)
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

GRANT ALL ON TABLE public.users TO session_manager;

GRANT ALL ON TABLE public.users TO postgres;

