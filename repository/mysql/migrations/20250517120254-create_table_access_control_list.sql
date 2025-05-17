-- +migrate Up
CREATE TABLE access_controls
(
    id            INT PRIMARY KEY AUTO_INCREMENT,
    functor_id    INT                   NOT NULL UNIQUE,
    functor_type  ENUM ('role', 'user') NOT NULL,
    permission_id INT                   NOT NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

-- +migrate Down
DROP TABLE access_controls;