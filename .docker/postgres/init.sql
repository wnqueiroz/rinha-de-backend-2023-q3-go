ALTER SYSTEM SET max_connections = 1000;
ALTER DATABASE rinha SET synchronous_commit = OFF;
ALTER SYSTEM SET shared_buffers TO "450MB";

CREATE TABLE public.pessoas (
	id UUID PRIMARY KEY NOT NULL,
	apelido VARCHAR(32) UNIQUE NOT NULL,
	nome VARCHAR(100) NOT NULL,
	nascimento DATE NOT NULL,
	stack TEXT NULL,
	search_trgm VARCHAR(1200) NOT NULL
);

CREATE extension IF NOT EXISTS pg_trgm;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_gin_pessoas_on_search_trgm ON pessoas USING gin (search_trgm gin_trgm_ops);

CREATE OR REPLACE FUNCTION notify_person_created() RETURNS TRIGGER as $notify_person_created$
BEGIN
  IF (TG_OP = 'INSERT') THEN
    PERFORM PG_NOTIFY(
      'person_created',
      JSON_BUILD_OBJECT(
        'id', NEW.id,
        'nome', NEW.nome,
        'apelido', NEW.apelido,
        'nascimento', NEW.nascimento,
        'stack', NEW.stack
      )::TEXT
    );
  END IF;
  RETURN NEW;
END;
$notify_person_created$ LANGUAGE plpgsql;

CREATE TRIGGER notify_person_created AFTER INSERT ON pessoas FOR EACH ROW EXECUTE PROCEDURE notify_person_created();