CREATE DATABASE loan_eligibilty WITH ENCODING='UTF8';
\c loan_user;
CREATE USER loan_user WITH PASSWORD '12345';
REVOKE ALL ON SCHEMA PUBLIC FROM PUBLIC;
REVOKE CONNECT ON DATABASE loan_eligibility FROM PUBLIC;
GRANT ALL PRIVILEGES ON DATABASE loan_eligibility TO loan_user;
GRANT ALL PRIVILEGES ON SCHEMA public TO loan_user;