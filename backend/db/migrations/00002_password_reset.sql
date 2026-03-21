-- +goose Up
CREATE TABLE password_reset_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_tbl(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tbl_token_hash ON password_reset_tbl(token_hash);
CREATE INDEX idx_password_reset_tbl_user_id ON password_reset_tbl(user_id);

-- +goose Down
DROP TABLE IF EXISTS password_reset_tbl;
