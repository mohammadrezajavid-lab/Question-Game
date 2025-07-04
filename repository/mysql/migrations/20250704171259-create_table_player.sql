-- +migrate Up
CREATE TABLE players
(
    id      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    game_id BIGINT UNSIGNED NOT NULL,
    score   INT UNSIGNED DEFAULT 0,

    UNIQUE KEY uq_user_game (user_id, game_id),

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,

    INDEX idx_user (user_id),
    INDEX idx_game (game_id)
);

CREATE TABLE player_answers
(
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    player_id   BIGINT UNSIGNED  NOT NULL,
    question_id BIGINT UNSIGNED  NOT NULL,
    choice      TINYINT UNSIGNED NOT NULL CHECK ( choice BETWEEN 1 AND 4),

    UNIQUE KEY uq_player_answer (player_id, question_id, choice),

    FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions (id),

    INDEX idx_player (player_id),
    INDEX idx_question (question_id)
);

-- +migrate Down
DROP TABLE players;
DROP TABLE player_answers;