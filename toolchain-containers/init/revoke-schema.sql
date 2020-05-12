DO
$do$
begin
    revoke all on schema public from public;
    revoke connect on database loan_eligibility from public;
    GRANT ALL PRIVILEGES ON DATABASE loan_eligibility TO loan_user;
    GRANT ALL PRIVILEGES ON SCHEMA public TO loan_user;
END
$do$;
