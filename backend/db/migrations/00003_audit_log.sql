-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_log_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action VARCHAR(50) NOT NULL,
    actor_id UUID,
    tenant_id UUID,
    ip_address VARCHAR(45),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_log_action ON audit_log_tbl(action);
CREATE INDEX idx_audit_log_actor ON audit_log_tbl(actor_id);
CREATE INDEX idx_audit_log_created ON audit_log_tbl(created_at);
CREATE INDEX idx_audit_log_tenant_created ON audit_log_tbl(tenant_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_log_tbl;
-- +goose StatementEnd
