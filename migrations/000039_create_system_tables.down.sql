DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS services;
ALTER TABLE users DROP FOREIGN KEY fk_users_department_id;
ALTER TABLE users DROP COLUMN department_id;
DROP TABLE IF EXISTS departments;
