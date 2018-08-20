# database

## create
```sql
CREATE DATABASE db_name
  WITH OWNER = postgres
       ENCODING = 'UTF8'

COMMENT ON DATABASE db_name
  IS 'xxx';
```
```sql
```

## alter
```
ALTER DATABASE db_name RENAME TO new_db_name;
ALTER DATABASE db_name OWNER  TO new_owner;
```

# table

## alter
```sql
ALTER TABLE tb_name RENAME TO new_tb_name;
ALTER TABLE tb_name ALTER COLUMN col_name TYPE col_type;
ALTER TABLE tb_name RENAME col_name TO new_col_name col_type;
ALTER TABLE tb_name ADD COLUMN col_name col_type;
```
