package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/mutation"
	"saas-template/internal/webserver/middleware"
	"strings"

	"github.com/google/uuid"
)

func Log(
	ctx context.Context,
	app *config.App,
	action string,
	actorID *uuid.UUID,
	tenantID *uuid.UUID,
	metadata map[string]string,
) {
	ip := extractIP(ctx)

	auditLog := model.AuditLogTbl{
		Action:    action,
		ActorID:   actorID,
		TenantID:  tenantID,
		IPAddress: ip,
	}

	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err == nil {
			s := string(b)
			auditLog.Metadata = &s
		}
	}

	err := mutation.CreateAuditLog(ctx, app.DB(), auditLog)
	if err != nil {
		slog.Default().
			ErrorContext(ctx, "failed to create audit log", "error", err, "action", action)
	}
}

func extractIP(ctx context.Context) *string {
	r := middleware.GetRequest(ctx)
	if r == nil {
		return nil
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		return &ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		addr := r.RemoteAddr
		return &addr
	}
	return &ip
}
