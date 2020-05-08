DO
$do$
begin

   IF NOT EXISTS (
      SELECT FROM pg_user  -- SELECT list can be empty for this
      WHERE  usename = 'loan_user') THEN
      create user loan_user with login password '${loan_db_pass}';
   END IF;

   GRANT ALL PRIVILEGES ON DATABASE loan_eligibility TO loan_user;

END
$do$;

-- create user 'user_db' with encrypted password 'password';
-- grant all privileges on database loan_eligibility to user_db;