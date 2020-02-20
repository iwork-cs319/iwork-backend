INSERT INTO workspaces (id, name)
VALUES
    (uuid_generate_v3(uuid_ns_url(), 'https://www.i.work/floors/1'), 'Main Floor');