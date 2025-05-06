-- +migrate Up
CREATE TABLE Users
(
    id              int primary key auto_increment,
    name            varchar(255) not null,
    phone_number    varchar(255) not null unique,
    hashed_password varchar(500) not null,
    created_at      datetime     not null default current_timestamp
);

-- +migrate Down
DROP TABLE Users;