CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT gen_random_uuid(),
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

INSERT INTO users (id, username, password, first_name, last_name) VALUES ('8411ef80-89d2-47e0-ae00-5140d00cf5a5', 'aaaa', '$2a$08$HKmdlpn7Kam0qXuRGwo88u0RHMY9Rt1atnQvmuPK8K8Gfb3xVd88m', 'Adrian', 'Adams'), ('18b4a8c3-b2e1-446b-8f69-af4aa1be677d', 'bbbb', '$2a$08$lFdX31fjc.Mmm3xUOe5siuVL4Mnm3e6k4Q7RGprCe./ZoPVW/.18W', 'Bob', 'Burns'), ('96aa5338-dd38-47ee-8ecf-2d529022670d', 'cccc', '$2a$08$P/w9f7Jizguk9kZaCz8wxe/HIPcKSgM7H5y2Y0lR4NLiXWhYeO9eK', 'Charlie', 'Churns'); 

INSERT INTO roles (name) VALUES ('admin'), ('user');

INSERT INTO user_roles (user_id, role_id) VALUES ('8411ef80-89d2-47e0-ae00-5140d00cf5a5', 1), ('18b4a8c3-b2e1-446b-8f69-af4aa1be677d', 2), ('96aa5338-dd38-47ee-8ecf-2d529022670d', 2);
