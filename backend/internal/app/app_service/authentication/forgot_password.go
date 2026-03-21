package authentication

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/mutation"
	"saas-template/internal/app/repository"
	"saas-template/internal/app/util_service/email"
	"time"
)

func ForgotPassword(ctx context.Context, app *config.App, emailAddr string) {
	user, err := repository.GetUserByEmail(ctx, app.DB(), emailAddr)
	if err != nil {
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		slog.Default().ErrorContext(ctx, "failed to generate reset token", "error", err)
		return
	}
	rawToken := hex.EncodeToString(tokenBytes)

	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(1 * time.Hour)

	resetRecord := model.PasswordResetTbl{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := mutation.CreatePasswordReset(ctx, app.DB(), resetRecord); err != nil {
		slog.Default().ErrorContext(ctx, "failed to create password reset record", "error", err)
		return
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", app.EnvVars().AppBaseUrl(), rawToken)

	if app.EnvVars().ResendApiKey() != "" {
		if err := email.SendPasswordResetEmail(
			app.Resend(),
			app.EnvVars().ResendFromEmail(),
			user.Email,
			resetURL,
		); err != nil {
			slog.Default().ErrorContext(ctx, "failed to send reset email", "error", err)
			return
		}
	} else {
		slog.Default().Info("password reset email skipped (no RESEND_API_KEY)", "url", resetURL)
	}
}
