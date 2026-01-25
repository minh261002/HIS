-- Rollback seed data
DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE username = 'admin');
DELETE FROM users WHERE username = 'admin';
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM roles;
