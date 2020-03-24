INSERT
INTO workspaces (id, floor_id, name)
VALUES (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-001'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-002'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-003'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-004'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-005'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-006'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'W-007');

INSERT
INTO workspace_assignee (id, user_id, workspace_id, start_time, end_time)
VALUES (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/assignee/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'), '2018-12-01T00:00:00Z', null),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/assignee/2'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'), '2018-12-01T00:00:00Z', null),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/assignee/3'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), '2018-12-01T00:00:00Z', null);