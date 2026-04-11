UPDATE schema_migrations SET version = 1, dirty = FALSE;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;