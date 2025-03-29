DROP DATABASE IF EXISTS micro_users_db;
CREATE DATABASE micro_users_db
WITH OWNER = postgres
     ENCODING = 'UTF8'
     TABLESPACE = pg_default
     CONNECTION LIMIT = -1;
ALTER DATABASE micro_users_db WITH CONNECTION LIMIT = -1;
