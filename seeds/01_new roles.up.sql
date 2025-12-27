INSERT INTO "roles"
(id, name, description)
VALUES 
(1, 'admin', 'Administrator role'),
(2, 'viewer', 'Viewer role'),
(3, 'user', 'Standart role')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('roles', 'id'), 
              (SELECT MAX(id) FROM "roles"));