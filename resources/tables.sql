create extension if not exists "uuid-ossp";
DROP TABLE IF EXISTS offerings;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS workspace_assignee;
DROP TABLE IF EXISTS workspaces;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS floors;

CREATE TABLE floors
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         TEXT NOT NULL,
    download_url TEXT NOT NULL,
    address      TEXT NOT NULL,
    deleted      BOOLEAN          DEFAULT FALSE
);

CREATE TABLE users
(
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       TEXT NOT NULL,
    department TEXT NOT NULL,
    email      TEXT             DEFAULT '',
    is_admin   BOOLEAN,
    deleted    BOOLEAN          DEFAULT FALSE
);

insert into users(id, name, department, email, is_admin)
VALUES ('decade00-0000-4000-a000-000000000000', 'Default User', 'N/A', 'N/A', false);

CREATE TABLE workspaces
(
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    floor_id uuid REFERENCES floors (id) NOT NULL,
    name     TEXT                        NOT NULL,
    details  TEXT             DEFAULT '',
    metadata JSON             DEFAULT '{}'::json,
    deleted  BOOLEAN          DEFAULT FALSE
);

CREATE TABLE bookings
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      uuid REFERENCES users (id)      NOT NULL,
    workspace_id uuid REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN          DEFAULT FALSE,
    start_time   TIMESTAMPTZ                     NOT NULL,
    end_time     TIMESTAMPTZ                     NOT NULL,
    created_by   uuid REFERENCES users (id)      NOT NULL
);

CREATE TABLE offerings
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      uuid REFERENCES users (id)      NOT NULL,
    workspace_id uuid REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN          DEFAULT FALSE,
    start_time   TIMESTAMPTZ                     NOT NULL,
    end_time     TIMESTAMPTZ,
    created_by   uuid REFERENCES users (id)      NOT NULL
);

CREATE TABLE workspace_assignee
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      uuid REFERENCES users (id)      NOT NULL,
    workspace_id uuid REFERENCES workspaces (id) NOT NULL,
    start_time   TIMESTAMPTZ                     NOT NULL,
    end_time     TIMESTAMPTZ
);