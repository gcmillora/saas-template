package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

func SendPasswordResetEmail(
	client *resend.Client,
	fromEmail string,
	toEmail string,
	resetURL string,
) error {
	html := fmt.Sprintf(`<h2>Password Reset</h2>
<p>You requested a password reset. Click the link below to set a new password:</p>
<p><a href="%s">Reset Password</a></p>
<p>This link expires in 1 hour. If you did not request this, ignore this email.</p>`, resetURL)

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: "Reset your password",
		Html:    html,
	}

	_, err := client.Emails.Send(params)
	return err
}
