DROP TABLE IF EXISTS offerings;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS workspaces;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS floors;

CREATE TABLE floors
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(32) NOT NULL,
    download_url TEXT        NOT NULL
);

CREATE TABLE users
(
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(32) NOT NULL,
    department VARCHAR(32) NOT NULL,
    is_admin   BOOLEAN
);

CREATE TABLE workspaces
(
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    floor_id uuid REFERENCES floors (id) NOT NULL,
    user_id  uuid REFERENCES users (id),
    name     VARCHAR(32)                        NOT NULL,
    locked   BOOLEAN          DEFAULT FALSE
);

CREATE TABLE bookings
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      uuid REFERENCES users (id)      NOT NULL,
    workspace_id uuid REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN          DEFAULT FALSE,
    start_time   TIMESTAMPTZ                            NOT NULL,
    end_time     TIMESTAMPTZ                            NOT NULL
);

CREATE TABLE offerings
(
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      uuid REFERENCES users (id)      NOT NULL,
    workspace_id uuid REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN          DEFAULT FALSE,
    start_time   TIMESTAMPTZ                            NOT NULL,
    end_time     TIMESTAMPTZ                            NOT NULL
);