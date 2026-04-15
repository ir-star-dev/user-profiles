INSERT INTO roles (role) VALUES ('admin'), ('user');

INSERT INTO users (name, email, password, role_id)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$BekeuI1csXF3zQXLRsgecu3W3zG75Bx82P0H5OajMz2ksb5.0IcdG',
    (SELECT id FROM roles WHERE role = 'admin')
)
ON CONFLICT (email) DO NOTHING;