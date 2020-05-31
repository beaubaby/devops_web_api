DO
$SUBENV_PG_FUNC
begin
   if not exists (
      select from PG_USER
      where  USENAME = 'loan_user') then
      CREATE USER loan_user with PASSWORD '${SUBENV_loan_db_pass}';
   end IF;

   ALTER USER loan_user WITH PASSWORD '${SUBENV_loan_db_pass}';
END
$SUBENV_PG_FUNC;