-- +migrate Up
CREATE TABLE games
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category   VARCHAR(32) NOT NULL,
    winner_id  BIGINT UNSIGNED DEFAULT NULL,
    start_time DATETIME    NOT NULL
);

-- +migrate Down
DROP TABLE games;