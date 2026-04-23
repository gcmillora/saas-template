package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"saas-template/config"
	"saas-template/internal/app/app_service/audit"
	"saas-template/internal/app/mutation"
	"saas-template/internal/app/repository"
	"saas-template/internal/app/util_service/password"

	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordResult struct {
	ValidationErrors []string
}

func ResetPassword(
	ctx context.Context,
	app *config.App,
	token string,
	newPassword string,
	confirmPassword string,
) (*ResetPasswordResult, error) {
	if newPassword != confirmPassword {
		return &ResetPasswordResult{ValidationErrors: []string{"Passwords do not match"}}, nil
	}

	if errs := password.ValidateComplexity(newPassword); len(errs) > 0 {
		return &ResetPasswordResult{ValidationErrors: errs}, nil
	}

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	resetRecord, err := repository.GetPasswordResetByTokenHash(ctx, app.DB(), tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired reset token")
	}

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	if err := mutation.UpdateUserPassword(
		ctx,
		app.DB(),
		resetRecord.UserID,
		string(bcryptHash),
	); err != nil {
		return nil, errors.New("failed to update password")
	}

	if err := mutation.InvalidateAllUserTokens(ctx, app.DB(), resetRecord.UserID); err != nil {
		slog.ErrorContext(ctx, "failed to invalidate reset tokens", "error", err, "user_id", resetRecord.UserID)
	}

	audit.Log(ctx, app, "password_reset_complete", &resetRecord.UserID, nil, nil)
	return nil, nil
}
