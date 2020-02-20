WITH main_floor_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1')),
     barry_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
     bruce_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce')),
     clark_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark')),
     diana_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'))
INSERT
INTO workspaces (id, floor_id, user_id, name, locked)
VALUES (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1'), main_floor_id, null, 'W-001', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2'), main_floor_id, null, 'W-002', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'), main_floor_id, clark_id, 'W-003', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'), main_floor_id, bruce_id, 'W-004', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), main_floor_id, diana_id, 'W-005', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6'), main_floor_id, null, 'W-006', false),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7'), main_floor_id, null, 'W-007', false);