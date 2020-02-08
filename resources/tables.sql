CREATE TABLE floors
(
    id   UUID PRIMARY KEY,
    name VARCHAR(10) NOT NULL
);

CREATE TABLE workspaces
(
    id       UUID PRIMARY KEY,
    floor_id UUID REFERENCES floors (id) NOT NULL,
    user_id  UUID,
    name     VARCHAR(10)                 NOT NULL,
    locked   BOOLEAN DEFAULT FALSE
);


CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(10) NOT NULL,
    department VARCHAR(10) NOT NULL,
    is_admin   BOOLEAN
);

CREATE TABLE bookings
(
    id           UUID PRIMARY KEY,
    user_id      UUID REFERENCES users (id)      NOT NULL,
    workspace_id UUID REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN DEFAULT FALSE,
    start_time   TIMESTAMPTZ                     NOT NULL,
    end_time     TIMESTAMPTZ                     NOT NULL
);

CREATE TABLE offerings
(
    id           UUID PRIMARY KEY,
    user_id      UUID REFERENCES users (id)      NOT NULL,
    workspace_id UUID REFERENCES workspaces (id) NOT NULL,
    cancelled    BOOLEAN DEFAULT FALSE,
    start_time   TIMESTAMPTZ                     NOT NULL,
    end_time     TIMESTAMPTZ                     NOT NULL
);