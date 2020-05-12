DO
$do$
begin

   if not exists (
      select from pg_user  -- SELECT list can be empty for this
      where  usename = 'loan_user') then
      create user loan_user with password '${loan_db_pass}';
   end IF;

END
$do$;