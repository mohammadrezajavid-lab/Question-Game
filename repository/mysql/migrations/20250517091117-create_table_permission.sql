-- +migrate Up
CREATE TABLE permissions
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    title      VARCHAR(191) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- +migrate Down
DROP TABLE permissions;