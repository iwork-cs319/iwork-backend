INSERT
INTO workspaces (id, floor_id, user_id, name, locked)
VALUES (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), null, 'W-001', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), null, 'W-002', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark'), 'W-003', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce'), 'W-004', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'), 'W-005', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), null, 'W-006', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), null, 'W-007', false);