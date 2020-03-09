INSERT
INTO bookings (id, user_id, workspace_id, cancelled, start_time, end_time, created_by)
VALUES (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1'), false, '2019-01-14T00:00:00Z',
        '2019-01-15T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/2'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1'), false, '2019-01-24T00:00:00Z',
        '2019-01-28T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/3'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2'), false, '2019-01-17T00:00:00Z',
        '2019-01-18T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/4'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2'), false, '2019-01-22T00:00:00Z',
        '2019-01-23T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/5'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'), false, '2019-01-17T00:00:00Z',
        '2019-01-17T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/6'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'), false, '2019-01-22T00:00:00Z',
        '2019-01-24T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), false, '2019-01-13T00:00:00Z',
        '2019-01-13T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/8'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), false, '2019-01-27T00:00:00Z',
        '2019-01-28T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/9'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6'), false, '2019-01-15T00:00:00Z',
        '2019-01-16T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/10'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6'), false, '2019-01-23T00:00:00Z',
        '2019-01-24T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/11'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6'), false, '2019-01-27T00:00:00Z',
        '2019-01-28T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/12'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7'), false, '2019-01-18T00:00:00Z',
        '2019-01-21T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/13'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7'), false, '2019-01-25T00:00:00Z',
        '2019-01-26T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'))
;

INSERT
INTO offerings (id, user_id, workspace_id, cancelled, start_time, end_time, created_by)
VALUES (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'), false, '2019-01-15T00:00:00Z',
        '2019-01-17T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/2'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3'), false, '2019-01-22T00:00:00Z',
        '2019-01-27T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/3'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'), false, '2019-01-16T00:00:00Z',
        '2019-01-19T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/4'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4'), false, '2019-01-22T00:00:00Z',
        '2019-01-24T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/5'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), false, '2019-01-13T00:00:00Z',
        '2019-01-14T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/6'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), false, '2019-01-22T00:00:00Z',
        '2019-01-24T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/7'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'),
        uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5'), false, '2019-01-27T00:00:00Z',
        '2019-01-28T11:59:59Z', uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry'))
;