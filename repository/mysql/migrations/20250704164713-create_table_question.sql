-- +migrate Up
CREATE TABLE questions
(
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    text              TEXT             NOT NULL,
    correct_answer_id BIGINT UNSIGNED,
    difficulty        TINYINT UNSIGNED NOT NULL CHECK ( difficulty BETWEEN 1 AND 3),
    category          VARCHAR(32)      NOT NULL,

    INDEX idx_category (category),
    INDEX idx_difficulty (difficulty)
);

CREATE TABLE possible_answers
(
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    text        TEXT             NOT NULL,
    question_id BIGINT UNSIGNED  NOT NULL,
    choice      TINYINT UNSIGNED NOT NULL CHECK ( choice BETWEEN 1 AND 4),

    FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE,
    UNIQUE (question_id, choice)
);

-- +migrate Down
DROP TABLE questions;
DROP TABLE possible_answers;
