DO
$do$
begin
    DROP schema IF EXISTS  public CASCADE;
    CREATE SCHEMA IF NOT EXISTS loan_schema;
END
$do$;