INSERT INTO users(name, department, is_admin)
VALUES ('Prayansh', 'IT', false)
RETURNING id;

INSERT INTO users(name, department, is_admin)
VALUES ('Mingyu', 'IT', false)
RETURNING id;

INSERT INTO users(name, department, is_admin)
VALUES ('Matt', 'IT', true)
RETURNING id;