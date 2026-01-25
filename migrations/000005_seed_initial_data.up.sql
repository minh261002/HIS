-- Seed initial roles
INSERT INTO roles (name, code, description, is_active, created_at, updated_at) VALUES
('Super Admin', 'SUPER_ADMIN', 'Full system access with all permissions', true, NOW(), NOW()),
('Admin', 'ADMIN', 'Administrative access for user and role management', true, NOW(), NOW()),
('Doctor', 'DOCTOR', 'Medical doctor role', true, NOW(), NOW()),
('Nurse', 'NURSE', 'Nursing staff role', true, NOW(), NOW()),
('Receptionist', 'RECEPTIONIST', 'Front desk and patient registration', true, NOW(), NOW());

-- Seed initial permissions
INSERT INTO permissions (name, code, description, module, created_at, updated_at) VALUES
-- User management permissions
('View Users', 'users.view', 'View user list and details', 'users', NOW(), NOW()),
('Create Users', 'users.create', 'Create new users', 'users', NOW(), NOW()),
('Update Users', 'users.update', 'Update user information', 'users', NOW(), NOW()),
('Delete Users', 'users.delete', 'Delete users', 'users', NOW(), NOW()),
('Manage Users', 'users.manage', 'Full user management access', 'users', NOW(), NOW()),

-- Role management permissions
('View Roles', 'roles.view', 'View role list and details', 'roles', NOW(), NOW()),
('Create Roles', 'roles.create', 'Create new roles', 'roles', NOW(), NOW()),
('Update Roles', 'roles.update', 'Update role information', 'roles', NOW(), NOW()),
('Delete Roles', 'roles.delete', 'Delete roles', 'roles', NOW(), NOW()),

-- Permission management permissions
('View Permissions', 'permissions.view', 'View permission list', 'permissions', NOW(), NOW()),
('Create Permissions', 'permissions.create', 'Create new permissions', 'permissions', NOW(), NOW()),
('Update Permissions', 'permissions.update', 'Update permissions', 'permissions', NOW(), NOW()),
('Delete Permissions', 'permissions.delete', 'Delete permissions', 'permissions', NOW(), NOW());

-- Assign all permissions to SUPER_ADMIN role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'SUPER_ADMIN';

-- Assign user and role management permissions to ADMIN role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'ADMIN'
AND p.code IN ('users.view', 'users.create', 'users.update', 'users.delete', 'users.manage', 
               'roles.view', 'roles.create', 'roles.update', 'roles.delete');

-- Create default admin user
-- Password: Admin@123 (bcrypt hash with cost 10)
INSERT INTO users (username, email, password_hash, full_name, phone_number, is_active, created_at, updated_at) VALUES
('admin', 'admin@hospital.local', '$2a$10$5k49ma4dtW86C1ymX8J1muObe0py8x1qxJgsBb0nhADzZV9ieUrUu', 'System Administrator', '', true, NOW(), NOW());

-- Assign SUPER_ADMIN role to admin user
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u
CROSS JOIN roles r
WHERE u.username = 'admin' AND r.code = 'SUPER_ADMIN';
