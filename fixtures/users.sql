INSERT INTO users (id, name, department, is_admin)
VALUES
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        'Barry Allen', 'R&D', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce'),
        'Bruce Wayne', 'Engineering', true),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark'),
        'Clark Kent', 'Marketing', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'),
        'Diana Prince', 'Operations', false);