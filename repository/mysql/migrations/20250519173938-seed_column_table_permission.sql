-- +migrate Up
INSERT INTO permissions(id, title)
VALUES (1, "user-add"),
       (2, "user-delete"),
       (3, "user-edit"),
       (4, "user-list");
-- +migrate Down
DELETE
FROM permissions
WHERE id IN (1, 2, 3, 4);