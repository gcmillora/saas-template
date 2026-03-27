-- +goose Up
ALTER TABLE user_tbl ADD CONSTRAINT chk_user_role CHECK (role IN ('admin', 'user'));

-- +goose Down
ALTER TABLE user_tbl DROP CONSTRAINT chk_user_role;
