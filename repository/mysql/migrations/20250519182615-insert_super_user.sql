-- +migrate Up
INSERT INTO users(id, name, phone_number, hashed_password, role)
VALUES (1, 'admin', '+989196551929', '$2a$14$P1bt5hURiXsnsxobjAHrweBpfr9MYG7Og4Pa5w0CzNpuTRMSxBdnO', 'admin');
-- +migrate Down
DELETE FROM users WHERE phone_number = '+989196551929';