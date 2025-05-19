-- +migrate Up
ALTER TABLE game_app_db.access_controls
    ADD CONSTRAINT unique_functor_permission UNIQUE (functor_id, functor_type, permission_id);

-- +migrate Down
ALTER TABLE game_app_db.access_controls DROP INDEX unique_functor_permission;