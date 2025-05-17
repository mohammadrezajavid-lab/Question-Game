-- +migrate Up
-- VARCHAR(191) for utf8mb4
CREATE TABLE users
(
    id              INT PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(191) NOT NULL,
    phone_number    VARCHAR(191) NOT NULL UNIQUE,
    hashed_password VARCHAR(191) NOT NULL,
    created_at      DATETIME DEFAULT current_timestamp
);

-- +migrate Down
DROP TABLE users;