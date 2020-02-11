DROP TABLE IF EXISTS offerings;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS workspaces;
DROP TABLE IF EXISTS floors;

CREATE TABLE floors
(
    id   VARCHAR(16) PRIMARY KEY,
    name VARCHAR(32) NOT NULL
);

CREATE TABLE workspaces
(
    id       VARCHAR(16) PRIMARY KEY,
    floor_id VARCHAR(16) REFERENCES floors (id) NOT NULL,
    user_id  VARCHAR(16) REFERENCES users (id),
    name     VARCHAR(32)                        NOT NULL,
    locked   BOOLEAN DEFAULT FALSE
);


CREATE TABLE users
(
    id         VARCHAR(16) PRIMARY KEY,
    name       VARCHAR(32) NOT NULL,
    department VARCHAR(32) NOT NULL,
    is_admin   BOOLEAN
);

CREATE TABLE bookings
(
    id           VARCHAR(16) PRIMARY KEY,
    user_id      VARCHAR(16) REFERENCES users (id)      NOT NULL,
    workspace_id VARCHAR(16) REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN DEFAULT FALSE,
    start_time   TIMESTAMPTZ                            NOT NULL,
    end_time     TIMESTAMPTZ                            NOT NULL
);

CREATE TABLE offerings
(
    id           VARCHAR(16) PRIMARY KEY,
    user_id      VARCHAR(16) REFERENCES users (id)      NOT NULL,
    workspace_id VARCHAR(16) REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN DEFAULT FALSE,
    start_time   TIMESTAMPTZ                            NOT NULL,
    end_time     TIMESTAMPTZ                            NOT NULL
);