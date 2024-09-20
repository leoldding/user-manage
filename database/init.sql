CREATE TABLE IF NOT EXISTS users (
    id  uuid DEFAULT gen_random_uuid(),
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS roles (
    id serial,
    name VARCHAR NOT NULL,
    PRIMARY KEY(id) 
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id uuid,
    role_id integer,
    PRIMARY KEY(user_id, role_id),
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id),
    CONSTRAINT fk_role_id
        FOREIGN KEY(role_id)
            REFERENCES roles(id)
);
