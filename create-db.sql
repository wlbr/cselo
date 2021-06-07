-- psql postgres -f create-db.sql

DROP DATABASE IF EXISTS cselo;

DROP ROLE IF EXISTS cseloapp;

CREATE ROLE cseloapp LOGIN CREATEDB;

CREATE DATABASE cselo
  WITH OWNER = cseloapp
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'de_DE.UTF-8'
       LC_CTYPE = 'de_DE.UTF-8'
       CONNECTION LIMIT = -1;
GRANT ALL ON DATABASE cselo TO cseloapp;
REVOKE ALL ON DATABASE cselo FROM public;
