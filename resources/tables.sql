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
    address      TEXT NOT NULL
);

CREATE TABLE users
(
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       TEXT NOT NULL,
    department TEXT NOT NULL,
    email      TEXT             DEFAULT '',
    is_admin   BOOLEAN
);

CREATE TABLE workspaces
(
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    floor_id uuid REFERENCES floors (id) NOT NULL,
    name     TEXT                        NOT NULL,
    locked   BOOLEAN          DEFAULT FALSE,
    details  TEXT             DEFAULT '',
    metadata JSON             DEFAULT '{}'::json
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
    end_time     TIMESTAMPTZ                     NOT NULL,
    created_by   uuid REFERENCES users (id)      NOT NULL
);

CREATE TABLE workspace_assignee
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      uuid REFERENCES users (id)      NOT NULL,
    workspace_id uuid REFERENCES workspaces (id) NOT NULL,
    start_time   TIMESTAMPTZ                     NOT NULL,
    end_time     TIMESTAMPTZ
)