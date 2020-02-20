WITH w_1 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1')),
     w_2 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2')),
     w_3 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3')),
     w_4 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4')),
     w_5 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5')),
     w_6 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6')),
     w_7 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7')),
     barry_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
     bruce_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce')),
     clark_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark')),
     diana_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'))
INSERT
INTO bookings (id, user_id, workspace_id, cancelled, start_time, end_time)
VALUES
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/1'), barry_id, w_1, false, '2019-01-14 00:00:00-07', '2019-01-15 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/2'), barry_id, w_1, false, '2019-01-24 00:00:00-07', '2019-01-28 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/3'), barry_id, w_2, false, '2019-01-17 00:00:00-07', '2019-01-18 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/4'), barry_id, w_2, false, '2019-01-22 00:00:00-07', '2019-01-23 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/5'), barry_id, w_3, false, '2019-01-17 00:00:00-07', '2019-01-17 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/6'), barry_id, w_4, false, '2019-01-22 00:00:00-07', '2019-01-24 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_5, false, '2019-01-13 00:00:00-07', '2019-01-13 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_5, false, '2019-01-27 00:00:00-07', '2019-01-28 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_6, false, '2019-01-15 00:00:00-07', '2019-01-16 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_6, false, '2019-01-23 00:00:00-07', '2019-01-24 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_6, false, '2019-01-27 00:00:00-07', '2019-01-28 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_7, false, '2019-01-18 00:00:00-07', '2019-01-21 11:59:59-07'),
       (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/bookings/7'), barry_id, w_7, false, '2019-01-25 00:00:00-07', '2019-01-26 11:59:59-07')
;

WITH w_1 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/1')),
     w_2 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/2')),
     w_3 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/3')),
     w_4 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/4')),
     w_5 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/5')),
     w_6 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/6')),
     w_7 AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/workspaces/7')),
     barry_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/barry')),
     bruce_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/bruce')),
     clark_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/clark')),
     diana_id AS (SELECT uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/users/diana'))
INSERT
INTO bookings (id, user_id, workspace_id, cancelled, start_time, end_time)
VALUES
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), clark_id, w_3, false, '2019-01-15 00:00:00-07', '2019-01-17 11:59:59-07'),
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), clark_id, w_3, false, '2019-01-22 00:00:00-07', '2019-01-27 11:59:59-07'),
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), bruce_id, w_4, false, '2019-01-16 00:00:00-07', '2019-01-19 11:59:59-07'),
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), bruce_id, w_4, false, '2019-01-22 00:00:00-07', '2019-01-24 11:59:59-07'),
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), diana_id, w_5, false, '2019-01-13 00:00:00-07', '2019-01-14 11:59:59-07'),
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), diana_id, w_5, false, '2019-01-22 00:00:00-07', '2019-01-24 11:59:59-07'),
        (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/offerings/1'), diana_id, w_5, false, '2019-01-27 00:00:00-07', '2019-01-28 11:59:59-07'),
;