-- +migrate Up
CREATE TABLE game_question
(
    question_id BIGINT UNSIGNED NOT NULL,
    game_id     BIGINT UNSIGNED NOT NULL,

    FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,

    PRIMARY KEY (question_id, game_id)
);

-- +migrate Down
DROP TABLE game_question;