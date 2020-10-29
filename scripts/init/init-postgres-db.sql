CREATE
    DATABASE loan WITH ENCODING = 'UTF8';

\c loan;

CREATE
    SCHEMA IF NOT EXISTS "loan" AUTHORIZATION "postgres";

ALTER USER postgres
    SET
    search_path TO loan;
