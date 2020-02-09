INSERT INTO users(id, name, department, is_admin)
VALUES (1, 'Prayansh', 'IT', false)
RETURNING id;

INSERT INTO users(id, name, department, is_admin)
VALUES (2, 'Mingyu', 'IT', false)
RETURNING id;

INSERT INTO users(id, name, department, is_admin)
VALUES (3, 'Matt', 'IT', true)
RETURNING id;