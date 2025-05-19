-- +migrate Up
INSERT INTO access_controls(functor_id, functor_type, permission_id)
VALUES (2, 'role', 1),
       (2, 'role', 2),
       (2, 'role', 3),
       (2, 'role', 4);


-- +migrate Down
DELETE
FROM access_controls
WHERE functor_id = 2
  and functor_type = 'role'
  and permission_id IN (1, 2, 3, 4);