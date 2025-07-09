-- +migrate Up
CREATE TABLE games
(
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category     VARCHAR(32) NOT NULL,
    winner_id    BIGINT UNSIGNED DEFAULT NULL,
    start_time   DATETIME    NOT NULL,
    question_ids JSON            DEFAULT NULL
);

-- +migrate Down
DROP TABLE games;